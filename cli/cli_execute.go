package cli

import (
	"fmt"
	"github.com/atotto/clipboard"
	s "github.com/baopham/snip/snippet"
	"github.com/fatih/color"
	"github.com/urfave/cli"
	"os/exec"
	"strings"
)

func Execute(c *cli.Context) error {
	keyword := strings.TrimSpace(c.Args().First())
	placeholderMap := strings.TrimSpace(c.Args().Get(1))

	if keyword == "" {
		return MissingInfoError{Message: "Please specify your keyword"}
	}

	filePath, err := s.SnippetFile()

	if err != nil {
		return err
	}

	mapper := convertPlaceholderMap(placeholderMap)

	// TODO: allow to select from a list
	snippet, err := s.SearchExact(keyword, filePath)

	if err != nil {
		return err
	}

	if snippet == nil {
		return NotFoundSnippet{Keyword: keyword}
	}

	content := snippet.Build(mapper)

	message := fmt.Sprintf("`%s` has been saved to your clipboard", content)

	if c.Bool("output") {
		message = fmt.Sprintf("`%s` *output* has been saved to your clipboard", content)
		parts := strings.Fields(content)
		command := exec.Command(parts[0], parts[1:]...)
		output, err := command.Output()
		if err == nil {
			err = clipboard.WriteAll(string(output))
		}
	} else {
		err = clipboard.WriteAll(content)
	}

	if err != nil {
		return err
	}

	color.Green(message)

	return nil
}

func convertPlaceholderMap(placeholderMap string) map[string]string {
	mapper := make(map[string]string)

	placeholderPairs := strings.Fields(placeholderMap)

	for _, pair := range placeholderPairs {
		parts := strings.Split(pair, "=")
		mapper[parts[0]] = parts[1]
	}

	return mapper
}
