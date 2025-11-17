package allow

import (
	"github.com/docker/go-plugins-helpers/authorization"
	"github.com/jonasbroms/hbm/docker/allow/types"
)

func True(req authorization.Request, config *types.Config) *types.AllowResult {
	return &types.AllowResult{Allow: true}
}
