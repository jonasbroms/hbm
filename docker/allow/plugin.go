package allow

import (
	"fmt"
	"net/url"

	"github.com/docker/go-plugins-helpers/authorization"
	"github.com/kassisol/hbm/docker/allow/types"
	"github.com/kassisol/hbm/internal/image"
	policyobj "github.com/kassisol/hbm/object/policy"
	"github.com/kassisol/hbm/version"
	log "github.com/sirupsen/logrus"
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
		log.WithFields(log.Fields{
			"version": version.Version,
		}).Fatal(err)
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
