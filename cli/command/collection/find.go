package collection

import (
	"fmt"
	"log/slog"
	"os"

	collectionobj "github.com/jonasbroms/hbm/object/collection"
	"github.com/jonasbroms/hbm/pkg/adf"
	"github.com/juliengk/go-utils"
	"github.com/spf13/cobra"
)

func newFindCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "find [name]",
		Short: "Verify if collection exists in the whitelist",
		Long:  findDescription,
		Args:  cobra.ExactArgs(1),
		Run:   runFind,
	}

	return cmd
}

func runFind(cmd *cobra.Command, args []string) {
	defer utils.RecoverFunc()

	c, err := collectionobj.New("sqlite", adf.AppPath)
	if err != nil {
		slog.Error("Failed to initialize collection store", "error", err)
		os.Exit(1)
	}
	defer c.End()

	fmt.Println(c.Find(args[0]))
}

var findDescription = `
Verify if collection exists in the whitelist

`
