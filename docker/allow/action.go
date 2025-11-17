package allow

import (
	"fmt"

	"github.com/jonasbroms/hbm/docker/allow/types"
	policyobj "github.com/jonasbroms/hbm/object/policy"
	"github.com/jonasbroms/hbm/version"
	"github.com/juliengk/go-utils"
	log "github.com/sirupsen/logrus"
)

func Action(config *types.Config, action, cmd string) *types.AllowResult {
	defer utils.RecoverFunc()

	p, err := policyobj.New("sqlite", config.AppPath)
	if err != nil {
		log.WithFields(log.Fields{
			"version": version.Version,
		}).Fatal(err)
	}
	defer p.End()

	if !p.Validate(config.Username, "action", action, "") {
		return &types.AllowResult{
			Allow: false,
			Msg: map[string]string{
				"text": fmt.Sprintf("%s is not allowed", cmd),
			},
		}
	}

	return &types.AllowResult{Allow: true}
}
