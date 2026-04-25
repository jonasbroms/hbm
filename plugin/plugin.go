package plugin

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/docker/go-plugins-helpers/authorization"
	"github.com/jonasbroms/hbm/pkg/uri"
	"github.com/jonasbroms/hbm/storage"
	"github.com/jonasbroms/hbm/storage/driver"
)

type plugin struct {
	appPath            string
	skipEndpoints      []*regexp.Regexp
	mu                 sync.Mutex
	internalContainers map[string]bool
}

func stringInRegexpSlice(s string, regexps []*regexp.Regexp) bool {
	for _, re := range regexps {
		if re.MatchString(s) {
			return true
		}
	}
	return false
}

func NewPlugin(appPath string) (*plugin, error) {
	p := &plugin{
		appPath: appPath,
		skipEndpoints: []*regexp.Regexp{
			regexp.MustCompile(`^/_ping`),
			regexp.MustCompile(`^/distribution/(.+)/json`),
			regexp.MustCompile(`^/events`), // used by our own event-stream goroutine
		},
		internalContainers: make(map[string]bool),
	}

	go p.purgeStaleOwners()
	go p.watchContainerEvents()

	return p, nil
}

// dockerClient returns an HTTP client that connects to the Docker Unix socket.
// Callers that need streaming responses (event stream) should set Timeout to 0.
func (p *plugin) dockerClient(timeout time.Duration) http.Client {
	return http.Client{
		Timeout: timeout,
		Transport: &http.Transport{
			DialContext: func(ctx context.Context, _, _ string) (net.Conn, error) {
				return net.Dial("unix", "/var/run/docker.sock")
			},
		},
	}
}

func (p *plugin) waitForDocker() bool {
	httpc := p.dockerClient(2 * time.Second)
	for i := 0; i < 30; i++ {
		resp, err := httpc.Get("http://localhost/_ping")
		if err == nil {
			resp.Body.Close()
			if resp.StatusCode == http.StatusOK {
				return true
			}
		}
		time.Sleep(2 * time.Second)
	}
	return false
}

// purgeStaleOwners runs once at startup, checks every stored container ID
// against Docker's inspect endpoint, and removes records for containers that
// no longer exist. This cleans up records left behind by crashes, upgrades,
// or any event the plugin missed while it was not running.
func (p *plugin) purgeStaleOwners() {
	if !p.waitForDocker() {
		slog.Warn("Docker not available, skipping stale container owner purge")
		return
	}

	s, err := storage.NewDriver("sqlite", p.appPath)
	if err != nil {
		slog.Warn("Failed to open storage for stale owner purge", "error", err)
		return
	}
	defer s.End()

	ids := s.ListContainerOwnerIDs()
	if len(ids) == 0 {
		return
	}

	// Mark all IDs as internal so inspect calls below bypass ContainerOwner auth.
	p.mu.Lock()
	for _, id := range ids {
		p.internalContainers[id] = true
	}
	p.mu.Unlock()
	defer func() {
		p.mu.Lock()
		for _, id := range ids {
			delete(p.internalContainers, id)
		}
		p.mu.Unlock()
	}()

	httpc := p.dockerClient(5 * time.Second)
	var purged, backfilled int
	for _, id := range ids {
		resp, err := httpc.Get("http://localhost/containers/" + id + "/json")
		if err != nil {
			continue
		}
		if resp.StatusCode == http.StatusNotFound {
			resp.Body.Close()
			if err := s.RemoveContainerOwner(id); err != nil {
				slog.Warn("Failed to remove stale container owner", "container_id", id, "error", err)
				continue
			}
			purged++
			continue
		}
		// Container exists — backfill name for records migrated from pre-0.20 schema
		// where container_name was stored as a separate "name:xxx" row and is now empty.
		var result struct {
			Name string `json:"Name"`
		}
		if err := json.NewDecoder(resp.Body).Decode(&result); err == nil {
			name := strings.TrimPrefix(result.Name, "/")
			if name != "" {
				if err := s.BackfillContainerName(id, name); err != nil {
					slog.Warn("Failed to backfill container name", "container_id", id, "error", err)
				} else {
					backfilled++
				}
			}
		}
		resp.Body.Close()
	}

	if purged > 0 {
		slog.Info("Purged stale container ownership records", "count", purged)
	}
	if backfilled > 0 {
		slog.Info("Backfilled container names from pre-0.20 ownership records", "count", backfilled)
	}
}

// watchContainerEvents subscribes to the Docker event stream and removes
// ownership records when the daemon emits a "destroy" event for a container.
// This covers all removal paths: explicit docker rm, --rm auto-removal, and
// daemon-internal cleanup — all produce a destroy event regardless of origin.
// The goroutine reconnects automatically and uses ?since= to catch events
// missed during a reconnect window.
func (p *plugin) watchContainerEvents() {
	if !p.waitForDocker() {
		slog.Warn("Docker not available, event stream not started")
		return
	}

	s, err := storage.NewDriver("sqlite", p.appPath)
	if err != nil {
		slog.Error("Failed to open storage for event stream", "error", err)
		return
	}
	defer s.End()

	// Initialise since to now so that reconnects re-request all events that
	// occurred since this goroutine started, rather than only events after the
	// moment of reconnection.
	since := time.Now().Unix()
	for {
		if err := p.streamEvents(s, &since); err != nil && err != io.EOF {
			slog.Warn("Container event stream error, reconnecting", "error", err)
		}
		time.Sleep(2 * time.Second)
	}
}

func (p *plugin) streamEvents(s driver.Storager, since *int64) error {
	v := url.Values{}
	v.Set("filters", `{"type":["container"],"event":["destroy"]}`)
	if *since > 0 {
		v.Set("since", strconv.FormatInt(*since, 10))
	}

	// No timeout: this connection intentionally lives forever.
	httpc := p.dockerClient(0)
	resp, err := httpc.Get("http://localhost/events?" + v.Encode())
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("events endpoint returned %d: %s", resp.StatusCode, body)
	}

	slog.Info("Container event stream connected")

	dec := json.NewDecoder(resp.Body)
	for {
		var event struct {
			Action string `json:"Action"`
			Status string `json:"status"` // pre-1.22 Docker compat
			ID     string `json:"id"`
			Time   int64  `json:"time"`
			Actor  struct {
				ID string `json:"ID"`
			} `json:"Actor"`
		}
		if err := dec.Decode(&event); err != nil {
			return err
		}

		// Guard against zero timestamps from non-event JSON (e.g. error bodies)
		// overwriting a valid since value.
		if event.Time > 0 {
			*since = event.Time
		}

		action := event.Action
		if action == "" {
			action = event.Status
		}

		// Use Actor.ID as fallback when the legacy top-level id field is absent.
		id := event.ID
		if id == "" {
			id = event.Actor.ID
		}

		if action == "destroy" && id != "" {
			if err := s.RemoveContainerOwner(id); err != nil {
				slog.Warn("Failed to remove container owner on destroy", "container_id", id, "error", err)
				continue
			}
			slog.Info("Container ownership removed", "event_type", "container_ownership",
				"container_id", id, "trigger", "destroy_event")
		}
	}
}

func (p *plugin) AuthZReq(req authorization.Request) authorization.Response {
	uriinfo, err := uri.GetURIInfo(req)
	if err != nil {
		return authorization.Response{Err: err.Error()}
	}

	if req.RequestMethod == "OPTIONS" || stringInRegexpSlice(uriinfo.Path, p.skipEndpoints) {
		return authorization.Response{Allow: true}
	}

	// Allow the plugin's own container-inspect calls (name lookup and purge).
	if req.RequestMethod == "GET" {
		re := regexp.MustCompile(`^/containers/([^/]+)/json$`)
		if m := re.FindStringSubmatch(uriinfo.Path); m != nil {
			p.mu.Lock()
			internal := p.internalContainers[m[1]]
			p.mu.Unlock()
			if internal {
				return authorization.Response{Allow: true}
			}
		}
	}

	a, err := NewApi(&uriinfo, p.appPath)
	if err != nil {
		return authorization.Response{Err: err.Error()}
	}

	r := a.Allow(req)
	if r.Error != "" {
		return authorization.Response{Err: r.Error}
	}
	if !r.Allow {
		return authorization.Response{Msg: r.Msg["text"]}
	}

	return authorization.Response{Allow: true}
}

func (p *plugin) iscreatecontainer(req authorization.Request, u *url.URL) bool {
	if req.ResponseStatusCode != 201 {
		return false
	}
	avm := regexp.MustCompile(`^/v\d+\.\d+/containers/create`)
	return avm.MatchString(u.Path) || u.Path == "/containers/create"
}

// getContainerName inspects a newly-created container to get its name.
// The inspect call is marked internal so it bypasses HBM's own auth check.
func (p *plugin) getContainerName(containerID string) (string, error) {
	p.mu.Lock()
	p.internalContainers[containerID] = true
	p.mu.Unlock()
	defer func() {
		p.mu.Lock()
		delete(p.internalContainers, containerID)
		p.mu.Unlock()
	}()

	httpc := p.dockerClient(5 * time.Second)
	resp, err := httpc.Get(fmt.Sprintf("http://localhost/containers/%s/json", containerID))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("inspect returned %d", resp.StatusCode)
	}

	var result struct {
		Name string `json:"Name"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}

	return strings.TrimPrefix(result.Name, "/"), nil
}

func (p *plugin) setcontainerowner(cname string, req authorization.Request) {
	username := req.User
	if username == "" {
		username = "root"
	}

	var rjson struct {
		Id string
	}
	if err := json.Unmarshal(req.ResponseBody, &rjson); err != nil {
		slog.Warn("Failed to parse container create response", "error", err)
		return
	}

	if cname == "" {
		var err error
		cname, err = p.getContainerName(rjson.Id)
		if err != nil {
			slog.Warn("Failed to get container name via inspect", "container_id", rjson.Id, "error", err)
		}
	}

	s, err := storage.NewDriver("sqlite", p.appPath)
	if err != nil {
		slog.Error("Failed to open storage for container ownership", "error", err)
		return
	}
	defer s.End()

	if err := s.SetContainerOwner(username, cname, rjson.Id); err != nil {
		slog.Warn("Failed to record container ownership", "user", username,
			"container_id", rjson.Id, "error", err)
		return
	}

	slog.Info("Container ownership recorded", "event_type", "container_ownership",
		"user", username, "container_name", cname, "container_id", rjson.Id)
}

func (p *plugin) AuthZRes(req authorization.Request) authorization.Response {
	u, err := url.Parse(req.RequestURI)
	if err != nil {
		return authorization.Response{Allow: true, Msg: err.Error()}
	}

	if p.iscreatecontainer(req, u) {
		p.setcontainerowner(u.Query().Get("name"), req)
	}

	return authorization.Response{Allow: true}
}
