package main

import (
	snippetCli "github.com/baopham/snip/cli"
	"github.com/fatih/color"
	"github.com/urfave/cli"
	"os"
)

func Action(fn func(c *cli.Context) error) func(c *cli.Context) error {
	actor := func(c *cli.Context) error {
		err := fn(c)

		if err != nil {
			color.Red(err.Error())
		}

		return err
	}

	return actor
}

func main() {
	app := cli.NewApp()
	app.Usage = "Save snippets: commands, texts, emoji, etc."
	app.Commands = []cli.Command{
		{
			Name: "snip",
		},
		{
			Name:    "save",
			Aliases: []string{"s"},
			Usage:   `snippets save "lsof -i :{p}"`,
			Action:  Action(snippetCli.Save),
		},
	}

	app.Run(os.Args)
}
