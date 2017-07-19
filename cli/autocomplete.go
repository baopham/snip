package cli

import (
	"encoding/csv"
	"fmt"
	"gopkg.in/baopham/snip.v2/snippet"
	"github.com/urfave/cli"
	"os"
)

func Autocomplete(c *cli.Context) {
	filePath, err := snippet.SnippetFile()

	if err != nil {
		return
	}

	file, err := os.Open(filePath)

	if err != nil {
		return
	}

	reader := csv.NewReader(file)

	for {
		row, err := reader.Read()

		if err != nil {
			return
		}

		fmt.Println(row[0])
	}
}
