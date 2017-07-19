package cli

import (
	s "gopkg.in/baopham/snip.v2/snippet"
	"github.com/olekukonko/tablewriter"
	"github.com/urfave/cli"
	"os"
)

func List(c *cli.Context) error {
	filePath, err := s.SnippetFile()

	if err != nil {
		return err
	}

	snippets, err := s.GetAll(filePath)

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
