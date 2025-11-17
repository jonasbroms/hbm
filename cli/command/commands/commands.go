package commands

import (
	"github.com/jonasbroms/hbm/cli/command/collection"
	"github.com/jonasbroms/hbm/cli/command/config"
	"github.com/jonasbroms/hbm/cli/command/group"
	"github.com/jonasbroms/hbm/cli/command/policy"
	"github.com/jonasbroms/hbm/cli/command/resource"
	"github.com/jonasbroms/hbm/cli/command/server"
	"github.com/jonasbroms/hbm/cli/command/system"
	"github.com/jonasbroms/hbm/cli/command/user"
	"github.com/spf13/cobra"
)

func AddCommands(cmd *cobra.Command) {
	cmd.AddCommand(
		collection.NewCommand(),
		config.NewCommand(),
		group.NewCommand(),
		policy.NewCommand(),
		resource.NewCommand(),
		user.NewCommand(),
		server.NewServerCommand(),
		system.NewInfoCommand(),
		system.NewInitCommand(),
		system.NewVersionCommand(),
	)
}
