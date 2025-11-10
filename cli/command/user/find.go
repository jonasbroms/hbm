package user

import (
	"fmt"
	"log/slog"
	"os"

	userobj "github.com/jonasbroms/hbm/object/user"
	"github.com/jonasbroms/hbm/pkg/adf"
	"github.com/juliengk/go-utils"
	"github.com/spf13/cobra"
)

func newFindCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "find [name]",
		Short: "Verify if user exists in the whitelist",
		Long:  findDescription,
		Args:  cobra.ExactArgs(1),
		Run:   runFind,
	}

	return cmd
}

func runFind(cmd *cobra.Command, args []string) {
	defer utils.RecoverFunc()

	u, err := userobj.New("sqlite", adf.AppPath)
	if err != nil {
		slog.Error("Failed to initialize user store", "error", err)
		os.Exit(1)
	}
	defer u.End()

	fmt.Println(u.Find(args[0]))
}

var findDescription = `
Verify if user exists in the whitelist

`
