package cli

import (
	s "github.com/baopham/snip/snippet"
	"github.com/olekukonko/tablewriter"
	"github.com/urfave/cli"
	"os"
	"strings"
)

func Search(c *cli.Context) error {
	keyword := strings.TrimSpace(c.Args().First())

	if keyword == "" {
		return MissingArgumentError{Message: "Please specify your keyword"}
	}

	filePath, err := s.SnippetFile()

	if err != nil {
		return err
	}

	snippets, err := s.Search(keyword, filePath)

	if err != nil {
		return err
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Keyword", "Content", "Description"})
	table.SetAlignment(tablewriter.ALIGN_LEFT)

	for _, snippet := range snippets {
		table.Append([]string{
			snippet.Keyword,
			snippet.Content,
			snippet.Description,
		})
	}

	table.Render()

	return nil
}
