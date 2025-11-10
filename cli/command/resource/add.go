package resource

import (
	"fmt"
	"log/slog"
	"os"

	resourcepkg "github.com/jonasbroms/hbm/docker/resource"
	resourceobj "github.com/jonasbroms/hbm/object/resource"
	"github.com/jonasbroms/hbm/pkg/adf"
	"github.com/juliengk/go-utils"
	"github.com/spf13/cobra"
)

var (
	resourceAddType   string
	resourceAddValue  string
	resourceAddOption []string
)

func newAddCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add [name]",
		Short: "Add resource to the whitelist",
		Long:  addDescription,
		Args:  cobra.ExactArgs(1),
		Run:   runAdd,
	}

	flags := cmd.Flags()
	flags.StringVarP(&resourceAddType, "type", "t", "action", fmt.Sprintf("Set resource type (%s)", resourcepkg.SupportedDrivers("|")))
	flags.StringVarP(&resourceAddValue, "value", "v", "", "Set resource value")
	flags.StringSliceVarP(&resourceAddOption, "option", "o", []string{}, "Specify options")

	return cmd
}

func runAdd(cmd *cobra.Command, args []string) {
	defer utils.RecoverFunc()

	r, err := resourceobj.New("sqlite", adf.AppPath)
	if err != nil {
		slog.Error("Failed to initialize resource store", "error", err)
		os.Exit(1)
	}
	defer r.End()

	if err := r.Add(args[0], resourceAddType, resourceAddValue, resourceAddOption); err != nil {
		slog.Error("Failed to add resource", "error", err)
		os.Exit(1)
	}
}

var addDescription = `
Add resource to the whitelist

`
