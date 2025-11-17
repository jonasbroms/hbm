package policy

import (
	"fmt"

	policyobj "github.com/jonasbroms/hbm/object/policy"
	"github.com/jonasbroms/hbm/pkg/adf"
	"github.com/juliengk/go-utils"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func newFindCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "find [name]",
		Short: "Verify if policy exists",
		Long:  "Verify if policy exists",
		Args:  cobra.ExactArgs(1),
		Run:   runFind,
	}

	return cmd
}

func runFind(cmd *cobra.Command, args []string) {
	defer utils.RecoverFunc()

	p, err := policyobj.New("sqlite", adf.AppPath)
	if err != nil {
		log.Fatal(err)
	}
	defer p.End()

	fmt.Println(p.Find(args[0]))
}

var findDescription = `
Verify if policy exists

`
