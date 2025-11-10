package policy

import (
	"log/slog"
	"os"

	policyobj "github.com/jonasbroms/hbm/object/policy"
	"github.com/jonasbroms/hbm/pkg/adf"
	"github.com/juliengk/go-utils"
	"github.com/spf13/cobra"
)

func newRemoveCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "rm [name]",
		Aliases: []string{"remove"},
		Short:   "Remove policy",
		Long:    removeDescription,
		Args:    cobra.ExactArgs(1),
		Run:     runRemove,
	}

	return cmd
}

func runRemove(cmd *cobra.Command, args []string) {
	defer utils.RecoverFunc()

	p, err := policyobj.New("sqlite", adf.AppPath)
	if err != nil {
		slog.Error("Failed to initialize policy store", "error", err)
		os.Exit(1)
	}
	defer p.End()

	if err := p.Remove(args[0]); err != nil {
		slog.Error("Failed to remove policy", "error", err)
		os.Exit(1)
	}
}

var removeDescription = `
Remove policy

`
