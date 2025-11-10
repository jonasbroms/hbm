package config

import (
	"fmt"
	"log/slog"
	"os"

	configobj "github.com/jonasbroms/hbm/object/config"
	"github.com/jonasbroms/hbm/pkg/adf"
	"github.com/juliengk/go-utils"
	"github.com/spf13/cobra"
)

func newGetCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "get [key]",
		Aliases: []string{"find"},
		Short:   "Get config option value",
		Long:    getDescription,
		Args:    cobra.ExactArgs(1),
		Run:     runGet,
	}

	return cmd
}

func runGet(cmd *cobra.Command, args []string) {
	defer utils.RecoverFunc()

	c, err := configobj.New("sqlite", adf.AppPath)
	if err != nil {
		slog.Error("Failed to initialize config store", "error", err)
		os.Exit(1)
	}
	defer c.End()

	result, err := c.Get(args[0])
	if err != nil {
		slog.Error("Failed to get config", "error", err)
		os.Exit(1)
	}

	fmt.Println(result)
}

var getDescription = `
Get config option value

`
