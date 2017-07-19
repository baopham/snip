package cli

import (
	"fmt"
	s "gopkg.in/baopham/snip.v2/snippet"
	"github.com/fatih/color"
	"github.com/urfave/cli"
	"strings"
)

func Remove(c *cli.Context) error {
	keyword := strings.TrimSpace(c.Args().First())

	if keyword == "" {
		return MissingInfoError{Message: "Please specify the snippet keyword"}
	}

	filePath, err := s.SnippetFile()

	if err != nil {
		return err
	}

	snippet, err := s.SearchExact(keyword, filePath)

	if err != nil {
		return err
	}

	if snippet == nil {
		return NotFoundSnippet{Keyword: keyword}
	}

	err = snippet.Remove(filePath)

	if err != nil {
		return err
	}

	color.Green(fmt.Sprintf("Snippet '%s' is removed", snippet.Keyword))

	return nil
}
