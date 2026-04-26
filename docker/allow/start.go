package allow

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"path/filepath"
	"strings"
	"time"

	"github.com/docker/go-plugins-helpers/authorization"
	"github.com/jonasbroms/hbm/docker/allow/types"
)

type containerInspectResult struct {
	HostConfig struct {
		Binds []string `json:"Binds"`
	} `json:"HostConfig"`
}

// inspectContainer fetches container config from Docker via the Unix socket.
// Since the call originates with no TLS user, api.go maps it to root (admin),
// bypassing ownership checks without needing the internalContainers map.
func inspectContainer(containerID string) (*containerInspectResult, error) {
	httpc := http.Client{
		Timeout: 5 * time.Second,
		Transport: &http.Transport{
			DialContext: func(ctx context.Context, _, _ string) (net.Conn, error) {
				return net.Dial("unix", "/var/run/docker.sock")
			},
		},
	}
	resp, err := httpc.Get("http://localhost/containers/" + containerID + "/json")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("inspect returned %d", resp.StatusCode)
	}
	var result containerInspectResult
	return &result, json.NewDecoder(resp.Body).Decode(&result)
}

// ContainerStart extends ContainerOwner with live re-validation of bind-mount
// paths to prevent TOCTOU symlink attacks. Between creation and start an attacker
// can replace an allowed directory with a symlink to a restricted location. This
// function resolves symlinks at start time and re-validates the real target
// against the current policy, closing that window.
func ContainerStart(req authorization.Request, config *types.Config) *types.AllowResult {
	r := ContainerOwner(req, config)
	if !r.Allow {
		return r
	}

	if config.DisableOwnershipCheck {
		return r
	}

	u, err := url.ParseRequestURI(req.RequestURI)
	if err != nil {
		return &types.AllowResult{Allow: false, Msg: map[string]string{
			"text": fmt.Sprintf("failed to parse request URI: %s", err),
		}}
	}

	// Path is either /vX.Y/containers/{id}/start or /containers/{id}/start.
	// ContainerOwner already validated the format; mirror its extraction.
	ts := strings.Trim(u.Path, "/")
	up := strings.Split(ts, "/")
	if len(up) < 3 {
		return &types.AllowResult{Allow: false, Msg: map[string]string{
			"text": fmt.Sprintf("invalid container start URI: %s", u.Path),
		}}
	}
	containerID := up[2]

	cfg, err := inspectContainer(containerID)
	if err != nil {
		// Fail closed: if we can't verify the mounts, deny.
		return &types.AllowResult{Allow: false, Msg: map[string]string{
			"text": fmt.Sprintf("failed to inspect container before start: %s", err),
		}}
	}

	for _, bind := range cfg.HostConfig.Binds {
		hostPath := strings.SplitN(bind, ":", 2)[0]
		if !strings.HasPrefix(hostPath, "/") {
			continue // named volume — no path to validate
		}

		resolved, err := filepath.EvalSymlinks(hostPath)
		if err != nil {
			// Path no longer exists — deny.
			return &types.AllowResult{Allow: false, Msg: map[string]string{
				"text":           fmt.Sprintf("volume path %s is no longer accessible", hostPath),
				"resource_type":  "volume",
				"resource_value": hostPath,
			}}
		}

		if resolved == hostPath {
			continue // unchanged since create, still valid
		}

		// Symlink introduced after creation — re-validate the real target.
		if allowed, reason := AllowVolume(resolved, config); !allowed {
			return &types.AllowResult{Allow: false, Msg: map[string]string{
				"text":           fmt.Sprintf("volume %s resolves to %s which is not allowed", hostPath, resolved),
				"resource_type":  "volume",
				"resource_value": hostPath,
				"denial_detail":  reason,
			}}
		}
	}

	return &types.AllowResult{Allow: true}
}
