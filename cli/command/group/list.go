package group

import (
	"fmt"
	"log/slog"
	"os"
	"strings"
	"text/tabwriter"

	groupobj "github.com/jonasbroms/hbm/object/group"
	"github.com/jonasbroms/hbm/pkg/adf"
	"github.com/juliengk/go-utils"
	"github.com/spf13/cobra"
)

var groupListFilter []string

func newListCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "ls",
		Aliases: []string{"list"},
		Short:   "List whitelisted groups",
		Long:    listDescription,
		Args:    cobra.NoArgs,
		Run:     runList,
	}

	flags := cmd.Flags()
	flags.StringSliceVarP(&groupListFilter, "filter", "f", []string{}, "Filter output based on conditions provided")

	return cmd
}

func runList(cmd *cobra.Command, args []string) {
	defer utils.RecoverFunc()

	g, err := groupobj.New("sqlite", adf.AppPath)
	if err != nil {
		slog.Error("Failed to initialize group store", "error", err)
		os.Exit(1)
	}
	defer g.End()

	filters := utils.ConvertSliceToMap("=", groupListFilter)

	groups, err := g.List(filters)
	if err != nil {
		slog.Error("Failed to list groups", "error", err)
		os.Exit(1)
	}

	if len(groups) > 0 {
		w := tabwriter.NewWriter(os.Stdout, 20, 1, 2, ' ', 0)
		fmt.Fprintln(w, "NAME\tUSERS")

		for group, users := range groups {
			if len(users) > 0 {
				fmt.Fprintf(w, "%s\t%s\n", group, strings.Join(users, ", "))
			} else {
				fmt.Fprintf(w, "%s\n", group)
			}
		}

		w.Flush()
	}
}

var listDescription = `
List whitelisted groups

`
