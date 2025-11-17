package allow

import (
	"net/url"
	"strings"

	"github.com/docker/go-plugins-helpers/authorization"
	"github.com/jonasbroms/hbm/docker/allow/types"
	policyobj "github.com/jonasbroms/hbm/object/policy"
	"github.com/jonasbroms/hbm/version"
	log "github.com/sirupsen/logrus"
)

func ContainerOwner(req authorization.Request, config *types.Config) *types.AllowResult {
	ar := &types.AllowResult{Allow: false}

	p, err := policyobj.New("sqlite", config.AppPath)
	if err != nil {
		log.WithFields(log.Fields{
			"version": version.Version,
		}).Fatal(err)
	}
	defer p.End()

	u, err := url.ParseRequestURI(req.RequestURI)
	if err != nil {
		return ar
	}

	ts := strings.Trim(u.Path, "/")
	up := strings.Split(ts, "/") // api version / type / id
	if len(up) < 3 {
		return ar
	}
	if up[1] != "containers" {
		return ar
	}

	ar.Allow = p.ValidateOwner(config.Username, "containers", up[2])

	return ar
}
