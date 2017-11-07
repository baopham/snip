package cli

import (
	s "github.com/baopham/snip/snippet"
	"github.com/urfave/cli"
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

func getPlaceholderMapper(args cli.Args) map[string]string {
	mapper := make(map[string]string)
	pair := args.Get(1)

	for i := 2; pair != ""; i++ {
		parts := strings.Split(pair, "=")
		mapper[strings.TrimSpace(parts[0])] = strings.TrimSpace(parts[1])
		pair = args.Get(i)
	}

	return mapper
}
