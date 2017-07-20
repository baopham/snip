package cli

import (
	"github.com/urfave/cli"
	s "gopkg.in/baopham/snip.v3/snippet"
	"strings"
)

func getSnippetContent(c *cli.Context) (string, error) {
	keyword := strings.TrimSpace(c.Args().First())

	if keyword == "" {
		return "", MissingInfoError{Message: "Please specify your keyword"}
	}

	filePath, err := s.SnippetFile()

	if err != nil {
		return "", err
	}

	mapper := getPlaceholderMapper(c.Args())

	snippet, err := s.SearchExact(keyword, filePath)

	if err != nil {
		return "", err
	}

	if snippet == nil {
		return "", NotFoundSnippet{Keyword: keyword}
	}

	return snippet.Build(mapper), nil
}
