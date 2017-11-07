package cli

import (
	"fmt"
	"github.com/atotto/clipboard"
	"github.com/fatih/color"
	"github.com/urfave/cli"
)

func Generate(c *cli.Context) error {
	content, err := getSnippetContent(c)

	if err != nil {
		return err
	}

	message := fmt.Sprintf("`%s` has been saved to your clipboard", content)
	err = clipboard.WriteAll(content)

	if err != nil {
		return err
	}

	color.Green(message)

	return nil
}
