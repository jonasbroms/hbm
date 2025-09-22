package allow

import (
	"fmt"
	"log/slog"
	"net"
	"os"
	"path/filepath"
	"strings"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/go-connections/nat"
	"github.com/docker/go-plugins-helpers/authorization"
	"github.com/jonasbroms/hbm/docker/allow/types"
	policyobj "github.com/jonasbroms/hbm/object/policy"
	objtypes "github.com/jonasbroms/hbm/object/types"
	"github.com/jonasbroms/hbm/version"
	"github.com/juliengk/go-mount"
	"github.com/juliengk/go-utils"
	"github.com/juliengk/go-utils/json"
)

func ContainerCreate(req authorization.Request, config *types.Config) *types.AllowResult {
	type ContainerCreateConfig struct {
		container.Config
		HostConfig container.HostConfig
	}

	cc := &ContainerCreateConfig{}

	if err := json.Decode(req.RequestBody, cc); err != nil {
		return &types.AllowResult{Allow: false, Error: err.Error()}
	}

	defer utils.RecoverFunc()

	p, err := policyobj.New("sqlite", config.AppPath)
	if err != nil {
		slog.Error("Failed to create policy object", "version", version.Version, "error", err)
		os.Exit(1)
	}
	defer p.End()

	if len(cc.HostConfig.Binds) > 0 {
		for _, b := range cc.HostConfig.Binds {
			// Docker volume binding format: [host-src:]container-dest[:<options>]
			// Split on : but need to handle edge cases
			if b == "" {
				return &types.AllowResult{
					Allow: false,
					Msg: map[string]string{
						"text":           "Empty volume binding not allowed",
						"resource_type":  "volume",
						"resource_value": b,
					},
				}
			}

			vol := strings.Split(b, ":")

			// Ensure we have at least the source path
			if len(vol) == 0 || vol[0] == "" {
				return &types.AllowResult{
					Allow: false,
					Msg: map[string]string{
						"text":           fmt.Sprintf("Invalid volume binding format - missing source path: %s", b),
						"resource_type":  "volume",
						"resource_value": b,
					},
				}
			}

			// Convert to absolute path and clean
			absPath, err := filepath.Abs(vol[0])
			if err != nil {
				return &types.AllowResult{
					Allow: false,
					Msg: map[string]string{
						"text":           fmt.Sprintf("Invalid volume path %s: %v", vol[0], err),
						"resource_type":  "volume",
						"resource_value": b,
					},
				}
			}

			if !AllowVolume(absPath, config) {
				return &types.AllowResult{
					Allow: false,
					Msg: map[string]string{
						"text":           fmt.Sprintf("Volume %s is not allowed to be mounted", absPath),
						"resource_type":  "volume",
						"resource_value": b,
					},
				}
			}
		}
	}

	if len(cc.HostConfig.LogConfig.Type) > 0 {
		if !p.Validate(config.Username, "logdriver", cc.HostConfig.LogConfig.Type, "") {
			return &types.AllowResult{
				Allow: false,
				Msg: map[string]string{
					"text":           fmt.Sprintf("Log driver %s is not allowed", cc.HostConfig.LogConfig.Type),
					"resource_type":  "logdriver",
					"resource_value": cc.HostConfig.LogConfig.Type,
				},
			}
		}
	}

	if len(cc.HostConfig.LogConfig.Config) > 0 {
		for k, v := range cc.HostConfig.LogConfig.Config {
			los := fmt.Sprintf("%s=%s", k, v)

			if !p.Validate(config.Username, "logopt", los, "") {
				return &types.AllowResult{
					Allow: false,
					Msg: map[string]string{
						"text":           fmt.Sprintf("Log option %s is not allowed", los),
						"resource_type":  "logopt",
						"resource_value": los,
					},
				}
			}
		}
	}

	if cc.HostConfig.NetworkMode == "host" {
		if !p.Validate(config.Username, "config", "container_create_param_net_host", "") {
			return &types.AllowResult{
				Allow: false,
				Msg: map[string]string{
					"text":           "--net=\"host\" param is not allowed",
					"resource_type":  "config",
					"resource_value": "container_create_param_net_host",
				},
			}
		}
	}

	if len(cc.HostConfig.PortBindings) > 0 {
		for cp, pbs := range cc.HostConfig.PortBindings {
			for _, pb := range pbs {
				spb := GetPortBindingString(&pb)

				cps := cp.Port()
				var fp string
				if spb != "" {
					fp = spb
				} else {
					fp = cps
				}
				if !p.Validate(config.Username, "port", fp, "") {
					return &types.AllowResult{
						Allow: false,
						Msg: map[string]string{
							"text":           fmt.Sprintf("Port %s is not allowed to be published", fp),
							"resource_type":  "port",
							"resource_value": fp,
						},
					}
				}
			}
		}
	}

	if len(cc.HostConfig.VolumeDriver) > 0 {
		if !p.Validate(config.Username, "volumedriver", cc.HostConfig.VolumeDriver, "") {
			return &types.AllowResult{
				Allow: false,
				Msg: map[string]string{
					"text":           fmt.Sprintf("Volume driver %s is not allowed", cc.HostConfig.VolumeDriver),
					"resource_type":  "volumedriver",
					"resource_value": cc.HostConfig.VolumeDriver,
				},
			}
		}
	}

	if len(cc.HostConfig.CapAdd) > 0 {
		for _, c := range cc.HostConfig.CapAdd {
			if !p.Validate(config.Username, "capability", c, "") {
				return &types.AllowResult{
					Allow: false,
					Msg: map[string]string{
						"text":           fmt.Sprintf("Capability %s is not allowed", c),
						"resource_type":  "capability",
						"resource_value": c,
					},
				}
			}
		}
	}

	if len(cc.HostConfig.DNS) > 0 {
		for _, dns := range cc.HostConfig.DNS {
			if !p.Validate(config.Username, "dns", dns, "") {
				return &types.AllowResult{
					Allow: false,
					Msg: map[string]string{
						"text":           fmt.Sprintf("DNS server %s is not allowed", dns),
						"resource_type":  "dns",
						"resource_value": dns,
					},
				}
			}
		}
	}

	if cc.HostConfig.IpcMode == "host" {
		if !p.Validate(config.Username, "config", "container_create_param_ipc_host", "") {
			return &types.AllowResult{
				Allow: false,
				Msg: map[string]string{
					"text":           "--ipc=\"host\" param is not allowed",
					"resource_type":  "config",
					"resource_value": "container_create_param_ipc_host",
				},
			}
		}
	}

	if cc.HostConfig.OomScoreAdj != 0 {
		if !p.Validate(config.Username, "config", "container_create_param_oom_score_adj", "") {
			return &types.AllowResult{
				Allow: false,
				Msg: map[string]string{
					"text":           "--oom-score-adj param is not allowed",
					"resource_type":  "config",
					"resource_value": "container_create_param_oom_score_adj",
				},
			}
		}
	}

	if cc.HostConfig.PidMode == "host" {
		if !p.Validate(config.Username, "config", "container_create_param_pid_host", "") {
			return &types.AllowResult{
				Allow: false,
				Msg: map[string]string{
					"text":           "--pid=\"host\" param is not allowed",
					"resource_type":  "config",
					"resource_value": "container_create_param_pid_host",
				},
			}
		}
	}

	if cc.HostConfig.Privileged {
		if !p.Validate(config.Username, "config", "container_create_param_privileged", "") {
			return &types.AllowResult{
				Allow: false,
				Msg: map[string]string{
					"text":           "--privileged param is not allowed",
					"resource_type":  "config",
					"resource_value": "container_create_param_privileged",
				},
			}
		}
	}

	if cc.HostConfig.PublishAllPorts {
		if !p.Validate(config.Username, "config", "container_create_param_publish_all", "") {
			return &types.AllowResult{
				Allow: false,
				Msg: map[string]string{
					"text":           "--publish-all param is not allowed",
					"resource_type":  "config",
					"resource_value": "container_create_param_publish_all",
				},
			}
		}
	}

	if len(cc.HostConfig.SecurityOpt) > 0 {
		if !p.Validate(config.Username, "config", "container_create_param_security_opt", "") {
			return &types.AllowResult{
				Allow: false,
				Msg: map[string]string{
					"text":           "--security-opt param is not allowed",
					"resource_type":  "config",
					"resource_value": "container_create_param_security_opt",
				},
			}
		}
	}

	if len(cc.HostConfig.Tmpfs) > 0 {
		if !p.Validate(config.Username, "config", "container_create_param_tmpfs", "") {
			return &types.AllowResult{
				Allow: false,
				Msg: map[string]string{
					"text":           "--tmpfs param is not allowed",
					"resource_type":  "config",
					"resource_value": "container_create_param_tmpfs",
				},
			}
		}
	}

	if cc.HostConfig.UTSMode == "host" {
		if !p.Validate(config.Username, "config", "container_create_param_uts_host", "") {
			return &types.AllowResult{
				Allow: false,
				Msg: map[string]string{
					"text":           "--uts=\"host\" param is not allowed",
					"resource_type":  "config",
					"resource_value": "container_create_param_uts_host",
				},
			}
		}
	}

	if cc.HostConfig.UsernsMode == "host" {
		if !p.Validate(config.Username, "config", "container_create_param_userns_host", "") {
			return &types.AllowResult{
				Allow: false,
				Msg: map[string]string{
					"text":           "--userns=\"host\" param is not allowed",
					"resource_type":  "config",
					"resource_value": "container_create_param_userns_host",
				},
			}
		}
	}

	if len(cc.HostConfig.Sysctls) > 0 {
		if !p.Validate(config.Username, "config", "container_create_param_sysctl", "") {
			return &types.AllowResult{
				Allow: false,
				Msg: map[string]string{
					"text":           "--sysctl param is not allowed",
					"resource_type":  "config",
					"resource_value": "container_create_param_sysctl",
				},
			}
		}
	}

	if len(cc.HostConfig.Runtime) > 0 {
		if !p.Validate(config.Username, "runtime", cc.HostConfig.Runtime, "") {
			return &types.AllowResult{
				Allow: false,
				Msg: map[string]string{
					"text":           fmt.Sprintf("Runtime %s is not allowed", cc.HostConfig.Runtime),
					"resource_type":  "runtime",
					"resource_value": cc.HostConfig.Runtime,
				},
			}
		}
	}

	if len(cc.HostConfig.Devices) > 0 {
		for _, dev := range cc.HostConfig.Devices {
			if !p.Validate(config.Username, "device", dev.PathOnHost, "") {
				return &types.AllowResult{
					Allow: false,
					Msg: map[string]string{
						"text":           fmt.Sprintf("Device %s is not allowed to be exported", dev.PathOnHost),
						"resource_type":  "device",
						"resource_value": dev.PathOnHost,
					},
				}
			}
		}
	}

	if cc.HostConfig.OomKillDisable != nil && *cc.HostConfig.OomKillDisable {
		if !p.Validate(config.Username, "config", "container_create_param_oom_kill_disable", "") {
			return &types.AllowResult{
				Allow: false,
				Msg: map[string]string{
					"text":           "--oom-kill-disable param is not allowed",
					"resource_type":  "config",
					"resource_value": "container_create_param_oom_kill_disable",
				},
			}
		}
	}

	if len(cc.HostConfig.Mounts) > 0 {
		for _, mount := range cc.HostConfig.Mounts {
			if mount.Type == "bind" {
				// Validate mount source is not empty or whitespace
				if len(mount.Source) == 0 || strings.TrimSpace(mount.Source) == "" {
					return &types.AllowResult{
						Allow: false,
						Msg: map[string]string{
							"text":           "Empty or whitespace-only mount source not allowed",
							"resource_type":  "volume",
							"resource_value": mount.Source,
						},
					}
				}

				// Convert to absolute path and clean
				absPath, err := filepath.Abs(mount.Source)
				if err != nil {
					return &types.AllowResult{
						Allow: false,
						Msg: map[string]string{
							"text":           fmt.Sprintf("Invalid mount source path: %s", mount.Source),
							"resource_type":  "volume",
							"resource_value": mount.Source,
						},
					}
				}

				if !AllowVolume(absPath, config) {
					return &types.AllowResult{
						Allow: false,
						Msg: map[string]string{
							"text":           fmt.Sprintf("Volume %s is not allowed to be mounted", absPath),
							"resource_type":  "volume",
							"resource_value": mount.Source,
						},
					}
				}
			}

			if mount.Type == "tmpfs" {
				if !p.Validate(config.Username, "config", "container_create_param_tmpfs", "") {
					return &types.AllowResult{
						Allow: false,
						Msg: map[string]string{
							"text":           "--tmpfs param is not allowed",
							"resource_type":  "config",
							"resource_value": "container_create_param_tmpfs",
						},
					}
				}
			}
		}
	}

	if len(cc.User) > 0 {
		if cc.Config.User == "root" && !p.Validate(config.Username, "config", "container_create_param_user_root", "") {
			return &types.AllowResult{
				Allow: false,
				Msg: map[string]string{
					"text":           "Running as user \"root\" is not allowed. Please use --user=\"someuser\" param.",
					"resource_type":  "config",
					"resource_value": "container_create_param_user_root",
				},
			}
		}
	}

	if !AllowImage(cc.Image, config) {
		return &types.AllowResult{
			Allow: false,
			Msg: map[string]string{
				"text":           fmt.Sprintf("Image %s is not allowed", cc.Image),
				"resource_type":  "image",
				"resource_value": cc.Image,
			},
		}
	}

	return &types.AllowResult{Allow: true}
}

func ipisany(ipstr string) bool {
	ip := net.ParseIP(ipstr)
	return ip.IsUnspecified()
}

func GetPortBindingString(pb *nat.PortBinding) string {
	result := pb.HostPort

	if len(pb.HostIP) > 0 && !ipisany(pb.HostIP) {
		result = fmt.Sprintf("%s:%s", pb.HostIP, pb.HostPort)
	}

	return result
}

// AllowVolume checks if a volume path is allowed to be mounted based on policy
func AllowVolume(vol string, config *types.Config) bool {
	defer utils.RecoverFunc()

	// Resolve symlinks to prevent symlink-based directory traversal attacks
	resolvedVol, err := filepath.EvalSymlinks(vol)
	if err != nil {
		slog.Warn("Failed to resolve symlinks for volume path", "path", vol, "error", err)
		return false
	}

	p, err := policyobj.New("sqlite", config.AppPath)
	if err != nil {
		slog.Error("Failed to create policy object", "version", version.Version, "error", err)
		os.Exit(1)
	}
	defer p.End()

	// Check exact match
	if checkVolumePermissionWithPolicy(p, resolvedVol, resolvedVol, false, config) {
		return true
	}

	// Check recursive parent permissions
	parts := strings.Split(strings.TrimPrefix(resolvedVol, "/"), "/")
	currentPath := "/"

	for _, part := range parts {
		if part == "" {
			continue
		}
		currentPath = filepath.Join(currentPath, part)

		if checkVolumePermissionWithPolicy(p, resolvedVol, currentPath, true, config) {
			return true
		}
	}

	return false
}

// checkVolumePermissionWithPolicy checks the permissions against the policies
func checkVolumePermissionWithPolicy(p policyobj.Policy, vol string, pathToCheck string, recursive bool, config *types.Config) bool {
	vo := objtypes.VolumeOptions{
		Recursive: recursive,
		NoSuid:    HasNoSuidFlag(vol),
	}

	jsonVO := json.Encode(vo)
	opts := strings.TrimSpace(jsonVO.String())

	return p.Validate(config.Username, "volume", pathToCheck, opts)
}

// HasNoSuidFlag checks if path is mounted with nosuid flag
func HasNoSuidFlag(vol string) bool {
	result := false

	entries, err := mount.New()
	if err != nil {
		return false
	}

	entry, err := entries.Find(vol)
	if err == nil {
		result = entry.FindOption("nosuid")
	}

	return result
}
