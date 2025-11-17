package allow

import (
	"fmt"
	"net/url"
	"path"
	"strings"

	"github.com/docker/go-plugins-helpers/authorization"
	"github.com/juliengk/go-utils"
	"github.com/juliengk/go-utils/json"
	"github.com/kassisol/hbm/docker/allow/types"
	"github.com/kassisol/hbm/internal/image"
	policyobj "github.com/kassisol/hbm/object/policy"
	objtypes "github.com/kassisol/hbm/object/types"
	"github.com/kassisol/hbm/version"
	log "github.com/sirupsen/logrus"
)

func ImageCreate(req authorization.Request, config *types.Config) *types.AllowResult {
	u, err := url.ParseRequestURI(req.RequestURI)
	if err != nil {
		return &types.AllowResult{
			Allow: false,
			Msg: map[string]string{
				"text": "Could not parse URL query",
			},
		}
	}

	params := u.Query()

	if v, ok := params["fromImage"]; ok {
		if !AllowImage(v[0], config) {
			return &types.AllowResult{
				Allow: false,
				Msg: map[string]string{
					"text":           fmt.Sprintf("Image %s is not allowed to be pulled", v[0]),
					"resource_type":  "image",
					"resource_value": v[0],
				},
			}
		}
	}

	return &types.AllowResult{Allow: true}
}

func AllowImage(img string, config *types.Config) bool {
	defer utils.RecoverFunc()

	p, err := policyobj.New("sqlite", config.AppPath)
	if err != nil {
		log.WithFields(log.Fields{
			"version": version.Version,
		}).Fatal(err)
	}
	defer p.End()

	i := image.NewImage(img)

	if i.Official {
		if p.Validate(config.Username, "config", "image_create_official", "") {
			return true
		}
	}

	if len(i.Registry) > 0 {
		if p.Validate(config.Username, "registry", i.Registry, "") {
			return true
		}
	}

	if p.Validate(config.Username, "image", i.String(), "") {
		return true
	}

	io := &objtypes.ImageOptions{SubImages: true}
	jio := json.Encode(io)
	opts := strings.TrimSpace(jio.String())

	dir, _ := path.Split(i.String())
	if len(dir) > 0 && p.Validate(config.Username, "image", dir, opts) {
		return true
	}

	return false
}
