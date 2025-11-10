package user

import (
	"log/slog"
	"os"

	userobj "github.com/jonasbroms/hbm/object/user"
	"github.com/jonasbroms/hbm/pkg/adf"
	"github.com/juliengk/go-utils"
	"github.com/spf13/cobra"
)

func newAddCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add [name]",
		Short: "Add user to the whitelist",
		Long:  addDescription,
		Args:  cobra.ExactArgs(1),
		Run:   runAdd,
	}

	return cmd
}

func runAdd(cmd *cobra.Command, args []string) {
	defer utils.RecoverFunc()

	u, err := userobj.New("sqlite", adf.AppPath)
	if err != nil {
		slog.Error("Failed to initialize user store", "error", err)
		os.Exit(1)
	}
	defer u.End()

	if err := u.Add(args[0]); err != nil {
		slog.Error("Failed to add user", "error", err)
		os.Exit(1)
	}
}

var addDescription = `
Add user to the whitelist

`
