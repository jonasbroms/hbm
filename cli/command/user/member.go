package user

import (
	"log/slog"
	"os"

	userobj "github.com/jonasbroms/hbm/object/user"
	"github.com/jonasbroms/hbm/pkg/adf"
	"github.com/juliengk/go-utils"
	"github.com/spf13/cobra"
)

var (
	userMemberAdd    bool
	userMemberRemove bool
)

func newMemberCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "member [group] [user]",
		Short: "Manage user membership to group",
		Long:  memberDescription,
		Args:  cobra.ExactArgs(2),
		Run:   runMember,
	}

	flags := cmd.Flags()
	flags.BoolVarP(&userMemberAdd, "add", "a", false, "Add user to group")
	flags.BoolVarP(&userMemberRemove, "remove", "r", false, "Remove user from group")

	return cmd
}

func runMember(cmd *cobra.Command, args []string) {
	defer utils.RecoverFunc()

	u, err := userobj.New("sqlite", adf.AppPath)
	if err != nil {
		slog.Error("Failed to initialize user store", "error", err)
		os.Exit(1)
	}
	defer u.End()

	if userMemberAdd {
		if err := u.AddToGroup(args[1], args[0]); err != nil {
			slog.Error("Failed to add user to group", "error", err)
			os.Exit(1)
		}
	}
	if userMemberRemove {
		if err := u.RemoveFromGroup(args[1], args[0]); err != nil {
			slog.Error("Failed to remove user from group", "error", err)
			os.Exit(1)
		}
	}
}

var memberDescription = `
Manage user membership to group

`
