package main

import (
	"github.com/urfave/cli"
	"os"
)

func main() {
	app := cli.NewApp()
	app.Name = "Snippets"
	app.Usage = "Save snippets: commands, texts, emoji, etc."
	app.Action = func(c *cli.Context) error {

		return nil
	}

	app.Run(os.Args)
}
