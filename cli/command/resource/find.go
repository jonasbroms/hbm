package resource

import (
	"fmt"

	resourceobj "github.com/jonasbroms/hbm/object/resource"
	"github.com/jonasbroms/hbm/pkg/adf"
	"github.com/juliengk/go-utils"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func newFindCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "find [name]",
		Short: "Verify if resource exists in the whitelist",
		Long:  findDescription,
		Args:  cobra.ExactArgs(1),
		Run:   runFind,
	}

	return cmd
}

func runFind(cmd *cobra.Command, args []string) {
	defer utils.RecoverFunc()

	r, err := resourceobj.New("sqlite", adf.AppPath)
	if err != nil {
		log.Fatal(err)
	}
	defer r.End()

	fmt.Println(r.Find(args[0]))
}

var findDescription = `
Verify if resource exists in the whitelist

`
