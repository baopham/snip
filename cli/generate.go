package cli

import (
	"fmt"
	"github.com/atotto/clipboard"
	"github.com/fatih/color"
	"github.com/urfave/cli"
	"strings"
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

func getPlaceholderMapper(args cli.Args) map[string]string {
	mapper := make(map[string]string)
	i := 1
	pair := args.Get(i)

	for pair != "" {
		parts := strings.Split(pair, "=")
		mapper[strings.TrimSpace(parts[0])] = strings.TrimSpace(parts[1])
		pair = args.Get(i)
		i++
	}

	return mapper
}
