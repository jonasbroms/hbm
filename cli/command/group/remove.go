package group

import (
	"log/slog"
	"os"

	groupobj "github.com/jonasbroms/hbm/object/group"
	"github.com/jonasbroms/hbm/pkg/adf"
	"github.com/juliengk/go-utils"
	"github.com/spf13/cobra"
)

func newRemoveCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "rm [name]",
		Aliases: []string{"remove"},
		Short:   "Remove group from the whitelist",
		Long:    removeDescription,
		Args:    cobra.ExactArgs(1),
		Run:     runRemove,
	}

	return cmd
}

func runRemove(cmd *cobra.Command, args []string) {
	defer utils.RecoverFunc()

	g, err := groupobj.New("sqlite", adf.AppPath)
	if err != nil {
		slog.Error("Failed to initialize group store", "error", err)
		os.Exit(1)
	}
	defer g.End()

	if err := g.Remove(args[0]); err != nil {
		slog.Error("Failed to remove group", "error", err)
		os.Exit(1)
	}
}

var removeDescription = `
Remove a group. You cannot remove a group that is in use by a policy.

`
