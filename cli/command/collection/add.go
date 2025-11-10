package collection

import (
	"log/slog"
	"os"

	collectionobj "github.com/jonasbroms/hbm/object/collection"
	"github.com/jonasbroms/hbm/pkg/adf"
	"github.com/juliengk/go-utils"
	"github.com/spf13/cobra"
)

func newAddCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add [name]",
		Short: "Add collection to the whitelist",
		Long:  addDescription,
		Args:  cobra.ExactArgs(1),
		Run:   runAdd,
	}

	return cmd
}

func runAdd(cmd *cobra.Command, args []string) {
	defer utils.RecoverFunc()

	c, err := collectionobj.New("sqlite", adf.AppPath)
	if err != nil {
		slog.Error("Failed to initialize collection store", "error", err)
		os.Exit(1)
	}
	defer c.End()

	if err := c.Add(args[0]); err != nil {
		slog.Error("Failed to add collection", "error", err)
		os.Exit(1)
	}
}

var addDescription = `
Add collection to the whitelist

`
