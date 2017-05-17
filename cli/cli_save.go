package cli

import (
	. "github.com/baopham/snip/snippet"
	"github.com/fatih/color"
	"github.com/urfave/cli"
	"strings"
)

var Trim = strings.TrimSpace

func Save(c *cli.Context) error {
	args := c.Args()
	keyword, content, description := Trim(args.First()), Trim(args.Get(1)), Trim(args.Get(2))

	if keyword == "" {
		return MissingArgumentError{Message: "Please specify your keyword"}
	}

	if content == "" {
		return MissingArgumentError{Message: "Please specify your snippet"}
	}

	snippet := Snippet{
		Keyword:     keyword,
		Content:     content,
		Description: description,
	}

	filePath, err := SnippetFile()

	if err != nil {
		return err
	}

	if err := snippet.Save(filePath); err != nil {
		return err
	}

	color.Green("Saved: " + snippet.Content)

	return nil
}
