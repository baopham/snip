package main

import (
	snippetCli "gopkg.in/baopham/snip.v3/cli"
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
	app.Version = "3.0.0"
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
			Name:         "search",
			Aliases:      []string{"s"},
			Usage:        "search for snippets: snip search port",
			Action:       Action(snippetCli.Search),
			BashComplete: snippetCli.Autocomplete,
		},
		{
			Name:         "generate",
			Aliases:      []string{"g"},
			Usage:        "generate the snippet by keyword: snip g port p={9000}",
			Action:       Action(snippetCli.Generate),
			BashComplete: snippetCli.Autocomplete,
		},
		{
			Name:    "execute",
			Aliases: []string{"x"},
			Usage:   "execute the snippet by keyword: snip x port p={9000}",
			Action:  Action(snippetCli.Execute),
			Flags: []cli.Flag{
				cli.BoolFlag{
					Name:  "output, o",
					Usage: "execute the snippet and save the output to clipboard",
				},
				cli.BoolFlag{
					Name:  "force, f",
					Usage: "skip the prompt and force to execute",
				},
			},
			BashComplete: snippetCli.Autocomplete,
		},
		{
			Name:         "list",
			Aliases:      []string{"l"},
			Usage:        "list all saved snippets: snip list",
			Action:       Action(snippetCli.List),
			BashComplete: snippetCli.Autocomplete,
		},
		{
			Name:         "remove",
			Aliases:      []string{"r"},
			Usage:        "remove a saved snippet: snip remove port",
			Action:       Action(snippetCli.Remove),
			BashComplete: snippetCli.Autocomplete,
		},
	}

	app.Run(os.Args)
}
