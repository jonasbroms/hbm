package group

import (
	"fmt"
	"log/slog"
	"os"

	groupobj "github.com/jonasbroms/hbm/object/group"
	"github.com/jonasbroms/hbm/pkg/adf"
	"github.com/juliengk/go-utils"
	"github.com/spf13/cobra"
)

func newFindCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "find [name]",
		Short: "Verify if group exists in the whitelist",
		Long:  findDescription,
		Args:  cobra.ExactArgs(1),
		Run:   runFind,
	}

	return cmd
}

func runFind(cmd *cobra.Command, args []string) {
	defer utils.RecoverFunc()

	g, err := groupobj.New("sqlite", adf.AppPath)
	if err != nil {
		slog.Error("Failed to initialize group store", "error", err)
		os.Exit(1)
	}
	defer g.End()

	fmt.Println(g.Find(args[0]))
}

var findDescription = `
Verify if group exists in the whitelist

`
