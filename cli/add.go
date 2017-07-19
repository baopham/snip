package cli

import (
	s "gopkg.in/baopham/snip.v3/snippet"
	"github.com/fatih/color"
	"github.com/urfave/cli"
	"strings"
)

var Trim = strings.TrimSpace

func Add(c *cli.Context) error {
	keyword, content, description := Trim(c.String("keyword")),
		Trim(c.String("content")), Trim(c.String("desc"))

	if keyword == "" {
		return MissingInfoError{Message: "Please specify your keyword"}
	}

	if content == "" {
		return MissingInfoError{Message: "Please specify your snippet"}
	}

	snippet := s.Snippet{
		Keyword:     keyword,
		Content:     content,
		Description: description,
	}

	filePath, err := s.SnippetFile()

	if err != nil {
		return err
	}

	if err := snippet.Save(filePath); err != nil {
		return err
	}

	color.Green("Saved: " + snippet.Content)

	return nil
}
