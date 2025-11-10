package policy

import (
	"log/slog"
	"os"

	policyobj "github.com/jonasbroms/hbm/object/policy"
	"github.com/jonasbroms/hbm/pkg/adf"
	"github.com/juliengk/go-utils"
	"github.com/spf13/cobra"
)

var (
	policyAddGroup      string
	policyAddCollection string
)

func newAddCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add [name]",
		Short: "Add policy",
		Long:  addDescription,
		Args:  cobra.ExactArgs(1),
		Run:   runAdd,
	}

	flags := cmd.Flags()
	flags.StringVarP(&policyAddGroup, "group", "g", "", "Set group")
	flags.StringVarP(&policyAddCollection, "collection", "c", "", "Set collection")

	return cmd
}

func runAdd(cmd *cobra.Command, args []string) {
	defer utils.RecoverFunc()

	p, err := policyobj.New("sqlite", adf.AppPath)
	if err != nil {
		slog.Error("Failed to initialize policy store", "error", err)
		os.Exit(1)
	}
	defer p.End()

	if err := p.Add(args[0], policyAddGroup, policyAddCollection); err != nil {
		slog.Error("Failed to add policy", "error", err)
		os.Exit(1)
	}
}

var addDescription = `
Add policy

`
