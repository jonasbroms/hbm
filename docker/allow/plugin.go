package allow

import (
	"fmt"
	"log/slog"
	"net/url"
	"os"

	"github.com/docker/go-plugins-helpers/authorization"
	"github.com/jonasbroms/hbm/docker/allow/types"
	"github.com/jonasbroms/hbm/internal/image"
	policyobj "github.com/jonasbroms/hbm/object/policy"
	"github.com/jonasbroms/hbm/version"
)

func PluginPull(req authorization.Request, config *types.Config) *types.AllowResult {
	var names []string
	var valid bool

	u, err := url.ParseRequestURI(req.RequestURI)
	if err != nil {
		return &types.AllowResult{
			Allow: false,
			Msg: map[string]string{
				"text": fmt.Sprintf("Could not parse URL query %s", req.RequestURI),
			},
		}
	}

	params := u.Query()

	pluginName := params["remote"][0]

	i := image.NewImage(pluginName)

	if len(i.Registry) > 0 {
		i.Registry = ""

		names = append(names, i.String())
	}
	names = append(names, pluginName)

	p, err := policyobj.New("sqlite", config.AppPath)
	if err != nil {
		slog.Error("Failed to create policy object", "version", version.Version, "error", err)
		os.Exit(1)
	}
	defer p.End()

	for _, name := range names {
		if p.Validate(config.Username, "plugin", name, "") {
			valid = true
		}
	}

	if !valid {
		return &types.AllowResult{
			Allow: false,
			Msg: map[string]string{
				"text":           fmt.Sprintf("Plugin %s is not allowed to be installed", pluginName),
				"resource_type":  "plugin",
				"resource_value": pluginName,
			},
		}
	}

	return &types.AllowResult{Allow: true}
}
