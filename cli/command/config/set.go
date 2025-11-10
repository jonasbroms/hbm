package config

import (
	"log/slog"
	"os"

	configobj "github.com/jonasbroms/hbm/object/config"
	"github.com/jonasbroms/hbm/pkg/adf"
	"github.com/juliengk/go-utils"
	"github.com/spf13/cobra"
)

func newSetCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "set [key] [value]",
		Short: "Set HBM config option",
		Long:  setDescription,
		Args:  cobra.ExactArgs(2),
		Run:   runSet,
	}

	return cmd
}

func runSet(cmd *cobra.Command, args []string) {
	defer utils.RecoverFunc()

	c, err := configobj.New("sqlite", adf.AppPath)
	if err != nil {
		slog.Error("Failed to initialize config store", "error", err)
		os.Exit(1)
	}
	defer c.End()

	if err := c.Set(args[0], args[1]); err != nil {
		slog.Error("Failed to set config", "error", err)
		os.Exit(1)
	}
}

var setDescription = `
Set HBM config option

`
