package allow

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/jonasbroms/hbm/docker/allow/types"
	policyobj "github.com/jonasbroms/hbm/object/policy"
	"github.com/jonasbroms/hbm/version"
	"github.com/juliengk/go-utils"
)

func Action(config *types.Config, action, cmd string) *types.AllowResult {
	defer utils.RecoverFunc()

	p, err := policyobj.New("sqlite", config.AppPath)
	if err != nil {
		slog.Error("Failed to create policy object", "version", version.Version, "error", err)
		os.Exit(1)
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
