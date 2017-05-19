package cli

import (
	"fmt"
	"github.com/atotto/clipboard"
	s "github.com/baopham/snip/snippet"
	"github.com/fatih/color"
	"github.com/urfave/cli"
	"strings"
)

func Execute(c *cli.Context) error {
	keyword := strings.TrimSpace(c.Args().First())
	placeholderMap := strings.TrimSpace(c.Args().Get(1))

	if keyword == "" {
		return MissingInfoError{Message: "Please specify your keyword"}
	}

	if placeholderMap == "" {

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

	err = clipboard.WriteAll(content)

	if err != nil {
		return err
	}

	color.Green(fmt.Sprintf("`%s` has been saved to your clipboard", content))

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
