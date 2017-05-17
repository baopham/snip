package snippet

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/user"
	"path"
	"regexp"
	"strings"

	"github.com/baopham/snippets-cli/util"
)

// EOL end of line
const EOL = "\n"

// SnippetSeparator between 2 snippets
const SnippetSeparator = ">>>>>>"

// Snippet represents the snippet
type Snippet struct {
	Keyword     string
	Description string
	Content     string
}

// Save snippet
func (s *Snippet) Save(filePath string) error {
	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0600)

	if err != nil {
		return err
	}

	existingSnippet, err := SearchExact(s.Keyword, filePath)

	if err != nil {
		return err
	}

	if existingSnippet != nil {
		return SnippetAlreadyExistError{Keyword: existingSnippet.Keyword}
	}

	defer util.Check(file.Close)

	info, err := file.Stat()

	if err != nil {
		return err
	}

	content := s.String()

	if info.Size() > 0 {
		content = EOL + SnippetSeparator + EOL + content
	}

	if _, err = file.WriteString(content); err != nil {
		return err
	}

	return nil
}

// SearchExact exact search by keyword
func SearchExact(keyword string, filePath string) (*Snippet, error) {
	snippets, err := searchByKeyword(keyword, filePath, true)

	if err != nil {
		return nil, err
	}

	if len(snippets) == 0 {
		return nil, nil
	}

	return snippets[0], nil
}

// Search fuzzy search by keyword
func Search(keyword string, filePath string) ([]*Snippet, error) {
	return searchByKeyword(keyword, filePath, false)
}

func (s *Snippet) String() string {
	return fmt.Sprintf("%s|%s|%s", s.Keyword, s.Content, s.Description)
}

// Build snippet actual content using the given placeholders
func (s *Snippet) Build(placeholders map[string]string) string {
	content := s.Content

	for k, v := range placeholders {
		content = strings.Replace(content, fmt.Sprintf("{%s}", k), v, -1)
	}

	return content
}

// SnippetDir returns the default directory path to the saved snippets
func SnippetDir() (string, error) {
	currentUser, err := user.Current()

	if err != nil {
		return "", err
	}

	dir := path.Join(currentUser.HomeDir, ".snippets-cli")

	return dir, nil
}

// SnippetFile returns the default file path to the saved snippets
func SnippetFile() (string, error) {
	dir, err := SnippetDir()

	if err != nil {
		return "", err
	}

	filePath := path.Join(dir, "snippets")

	return filePath, err
}

func init() {
	dir, err := SnippetDir()

	if err != nil {
		panic(err)
	}

	if _, err = os.Stat(dir); os.IsNotExist(err) {
		err = os.MkdirAll(dir, 0755)

		if err != nil {
			panic(err)
		}
	}
}

func getScanner(file *os.File) *bufio.Scanner {
	const splitSeparator = EOL + SnippetSeparator + EOL

	trimSeparator := func(data []byte) []byte {
		return bytes.Trim(data, splitSeparator)
	}

	splitter := func(data []byte, atEOF bool) (advanced int, token []byte, err error) {
		if atEOF && len(data) == 0 {
			return 0, nil, nil
		}

		if i := strings.Index(string(data), splitSeparator); i >= 0 {
			return i + 1, trimSeparator(data[0:i]), nil
		}

		if atEOF {
			return len(data), trimSeparator(data), nil
		}

		return
	}

	scanner := bufio.NewScanner(file)

	scanner.Split(splitter)

	return scanner
}

func searchByKeyword(keyword string, filePath string, exact bool) ([]*Snippet, error) {
	var snippets []*Snippet

	file, err := os.Open(filePath)

	if err != nil {
		return snippets, err
	}

	scanner := getScanner(file)

	matcher := fuzzyMatcher

	if exact {
		matcher = exactMatcher
	}

	for line := 1; scanner.Scan(); line++ {
		lineContent := scanner.Text()

		if !matcher(keyword, lineContent) {
			continue
		}

		record := strings.Split(lineContent, "|")

		found := &Snippet{
			Keyword:     record[0],
			Content:     record[1],
			Description: record[2],
		}

		snippets = append(snippets, found)

		if exact {
			return snippets, nil
		}
	}

	return snippets, nil
}

func fuzzyMatcher(keyword string, content string) bool {
	return regexp.MustCompile(fmt.Sprintf(`^.*%s.*\|`, keyword)).MatchString(content)
}

func exactMatcher(keyword string, content string) bool {
	return regexp.MustCompile(fmt.Sprintf(`^%s\|`, keyword)).MatchString(content)
}
