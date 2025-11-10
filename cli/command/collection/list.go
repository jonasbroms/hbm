package collection

import (
	"fmt"
	"log/slog"
	"os"
	"strings"
	"text/tabwriter"

	collectionobj "github.com/jonasbroms/hbm/object/collection"
	"github.com/jonasbroms/hbm/pkg/adf"
	"github.com/juliengk/go-utils"
	"github.com/spf13/cobra"
)

var collectionListFilter []string

func newListCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "ls",
		Aliases: []string{"list"},
		Short:   "List whitelisted collections",
		Long:    listDescription,
		Args:    cobra.NoArgs,
		Run:     runList,
	}

	flags := cmd.Flags()
	flags.StringSliceVarP(&collectionListFilter, "filter", "f", []string{}, "Filter output based on conditions provided")

	return cmd
}

func runList(cmd *cobra.Command, args []string) {
	defer utils.RecoverFunc()

	c, err := collectionobj.New("sqlite", adf.AppPath)
	if err != nil {
		slog.Error("Failed to initialize collection store", "error", err)
		os.Exit(1)
	}
	defer c.End()

	filters := utils.ConvertSliceToMap("=", collectionListFilter)

	collections, err := c.List(filters)
	if err != nil {
		slog.Error("Failed to list collections", "error", err)
		os.Exit(1)
	}

	if len(collections) > 0 {
		w := tabwriter.NewWriter(os.Stdout, 20, 1, 2, ' ', 0)
		fmt.Fprintln(w, "NAME\tRESOURCES")

		for collection, resources := range collections {
			if len(resources) > 0 {
				fmt.Fprintf(w, "%s\t%s\n", collection, strings.Join(resources, ", "))
			} else {
				fmt.Fprintf(w, "%s\n", collection)
			}
		}

		w.Flush()
	}
}

var listDescription = `
List whitelisted collections

`
