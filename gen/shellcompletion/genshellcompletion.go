package main

import (
	"fmt"

	"os"

	"github.com/jonasbroms/hbm/cli/command"
	"github.com/jonasbroms/hbm/cli/command/commands"
)

func main() {
	scPath := "/tmp/hbm/shellcompletion"
	bashTarget := fmt.Sprintf("%s/bash", scPath)

	if err := os.MkdirAll(scPath, 0755); err != nil {
		fmt.Println(err)
	}

	cmd := command.NewHBMCommand()
	commands.AddCommands(cmd)
	cmd.DisableAutoGenTag = true

	if err := cmd.GenBashCompletionFile(bashTarget); err != nil {
		fmt.Println(err)
	}
}
