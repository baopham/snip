package cli

import (
	"fmt"
	"github.com/atotto/clipboard"
	"github.com/baopham/go-cliutil/cliutil"
	"github.com/fatih/color"
	"github.com/urfave/cli"
	"os"
	"os/exec"
)

func Execute(c *cli.Context) error {
	content, err := getSnippetContent(c)

	if err != nil {
		return err
	}

	if !c.Bool("force") {
		yes := cliutil.Prompt("Are you sure you want to execute: %s", content)

		if !yes {
			return nil
		}
	}

	cmd := exec.Command("bash", "-c", content)

	if c.Bool("output") {
		output, err := cmd.Output()
		if err != nil {
			return err
		}

		err = clipboard.WriteAll(string(output))

		if err != nil {
			return err
		}

		color.Green(fmt.Sprintf("`%s` *output* has been saved to your clipboard", content))
	} else {
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		return cmd.Run()
	}

	return nil
}
