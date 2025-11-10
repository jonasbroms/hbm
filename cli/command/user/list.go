package user

import (
	"fmt"
	"log/slog"
	"os"
	"strings"
	"text/tabwriter"

	userobj "github.com/jonasbroms/hbm/object/user"
	"github.com/jonasbroms/hbm/pkg/adf"
	"github.com/juliengk/go-utils"
	"github.com/spf13/cobra"
)

var userListFilter []string

func newListCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "ls",
		Aliases: []string{"list"},
		Short:   "List whitelisted users",
		Long:    listDescription,
		Args:    cobra.NoArgs,
		Run:     runList,
	}

	flags := cmd.Flags()
	flags.StringSliceVarP(&userListFilter, "filter", "f", []string{}, "Filter output based on conditions provided")

	return cmd
}

func runList(cmd *cobra.Command, args []string) {
	defer utils.RecoverFunc()

	u, err := userobj.New("sqlite", adf.AppPath)
	if err != nil {
		slog.Error("Failed to initialize user store", "error", err)
		os.Exit(1)
	}
	defer u.End()

	filters := utils.ConvertSliceToMap("=", userListFilter)

	users, err := u.List(filters)
	if err != nil {
		slog.Error("Failed to list users", "error", err)
		os.Exit(1)
	}

	if len(users) > 0 {
		w := tabwriter.NewWriter(os.Stdout, 20, 1, 2, ' ', 0)
		fmt.Fprintln(w, "NAME\tGROUPS")

		for user, groups := range users {
			if len(groups) > 0 {
				fmt.Fprintf(w, "%s\t%s\n", user, strings.Join(groups, ", "))
			} else {
				fmt.Fprintf(w, "%s\n", user)
			}
		}

		w.Flush()
	}
}

var listDescription = `
List whitelisted users

`
