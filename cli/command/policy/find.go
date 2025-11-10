package policy

import (
	"fmt"
	"log/slog"
	"os"

	policyobj "github.com/jonasbroms/hbm/object/policy"
	"github.com/jonasbroms/hbm/pkg/adf"
	"github.com/juliengk/go-utils"
	"github.com/spf13/cobra"
)

func newFindCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "find [name]",
		Short: "Verify if policy exists",
		Long:  "Verify if policy exists",
		Args:  cobra.ExactArgs(1),
		Run:   runFind,
	}

	return cmd
}

func runFind(cmd *cobra.Command, args []string) {
	defer utils.RecoverFunc()

	p, err := policyobj.New("sqlite", adf.AppPath)
	if err != nil {
		slog.Error("Failed to initialize policy store", "error", err)
		os.Exit(1)
	}
	defer p.End()

	fmt.Println(p.Find(args[0]))
}

var findDescription = `
Verify if policy exists

`
