package resource

import (
	"log/slog"
	"os"

	resourceobj "github.com/jonasbroms/hbm/object/resource"
	"github.com/jonasbroms/hbm/pkg/adf"
	"github.com/juliengk/go-utils"
	"github.com/spf13/cobra"
)

func newRemoveCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "rm [name]",
		Aliases: []string{"remove"},
		Short:   "Remove resource from the whitelist",
		Long:    removeDescription,
		Args:    cobra.ExactArgs(1),
		Run:     runRemove,
	}

	return cmd
}

func runRemove(cmd *cobra.Command, args []string) {
	defer utils.RecoverFunc()

	r, err := resourceobj.New("sqlite", adf.AppPath)
	if err != nil {
		slog.Error("Failed to initialize resource store", "error", err)
		os.Exit(1)
	}
	defer r.End()

	if err := r.Remove(args[0]); err != nil {
		slog.Error("Failed to remove resource", "error", err)
		os.Exit(1)
	}
}

var removeDescription = `
Remove resource from the whitelist

`
