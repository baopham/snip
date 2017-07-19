package cli

import (
	"fmt"
	s "gopkg.in/baopham/snip.v3/snippet"
	"github.com/urfave/cli"
	"strings"
)

func getSnippetContent(c *cli.Context) (string, error) {
	keyword := strings.TrimSpace(c.Args().First())

	if keyword == "" {
		return "", MissingInfoError{Message: "Please specify your keyword"}
	}

	filePath, err := s.SnippetFile()

	if err != nil {
		return "", err
	}

	mapper := getPlaceholderMapper(c.Args())

	snippet, err := s.SearchExact(keyword, filePath)

	if err != nil {
		return "", err
	}

	if snippet == nil {
		return "", NotFoundSnippet{Keyword: keyword}
	}

	return snippet.Build(mapper), nil
}

func prompt(format string, a ...interface{}) bool {
	if !strings.Contains(format, "[y/n]") {
		format += " [y/n] "
	}
	if len(a) == 0 {
		fmt.Print(format)
	} else {
		fmt.Printf(format, a...)
	}
	return handlePromptResponse()
}

func handlePromptResponse() bool {
	var response string
	_, err := fmt.Scanln(&response)

	if err != nil {
		return false
	}

	response = strings.TrimSpace(strings.ToLower(response))

	if response == "y" || response == "yes" {
		return true
	}

	return false
}
