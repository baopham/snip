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
	app.Version = "1.0.2"
	app.Usage = "Save snippets: commands, texts, emoji, etc."
	app.EnableBashCompletion = true
	app.Commands = []cli.Command{
		{
			Name:    "add",
			Aliases: []string{"a"},
			Usage:   `snip add -k="port" -c="lsof -i :{p}" -desc="List processes listening on a particular port"`,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "keyword, k",
					Usage: "keyword for the snippet",
				},
				cli.StringFlag{
					Name:  "content, c",
					Usage: "the snippet content",
				},
				cli.StringFlag{
					Name:  "description, desc",
					Usage: "the snippet description",
				},
			},
			Action: Action(snippetCli.Add),
		},
		{
			Name:    "search",
			Aliases: []string{"s"},
			Usage:   "search for snippets: snip search port",
			Action:  Action(snippetCli.Search),
		},
		{
			Name:    "execute",
			Aliases: []string{"x"},
			Usage:   "get the snippet by keyword: snip x port p={9000}",
			Flags: []cli.Flag{
				cli.BoolFlag{
					Name: "output, o",
					Usage: "execute right away and save the output to clipboard. " +
						"Without this option, the snippet will be saved to the clipboard instead",
				},
			},
			Action: Action(snippetCli.Execute),
		},
		{
			Name:    "list",
			Aliases: []string{"l"},
			Usage:   "list all saved snippets: snip list",
			Action:  Action(snippetCli.List),
		},
		{
			Name:    "remove",
			Aliases: []string{"r"},
			Usage:   "remove a saved snippet: snip remove port",
			Action:  Action(snippetCli.Remove),
		},
	}

	app.Run(os.Args)
}
