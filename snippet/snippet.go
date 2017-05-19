package snippet

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"os/user"
	"path"
	"regexp"
	"strings"

	"github.com/baopham/snip/util"
)

// EOL end of line
const EOL = "\n"

// SnippetSeparator between 2 snippets
const SnippetSeparator = ">>>>>>"

const FileMode = 0600

const (
	SEARCH_EXACT     SearchCode = 1
	SEARCH_FUZZY     SearchCode = 2
	SEARCH_MATCH_ANY SearchCode = 3
)

type SearchCode int

// Snippet represents the snippet
type Snippet struct {
	Keyword     string
	Description string
	Content     string
}

// Save snippet
func (s *Snippet) Save(filePath string) error {
	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, FileMode)

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

// Remove a saved snippet
func (s *Snippet) Remove(filePath string) error {
	b, err := ioutil.ReadFile(filePath)

	if err != nil {
		return err
	}

	content := string(b)

	re := regexp.MustCompile(fmt.Sprintf("(?m)(%s)*%s$%s*", SnippetSeparator+EOL, regexp.QuoteMeta(s.String()), EOL))

	newContent := re.ReplaceAllString(content, "")

	lines := strings.Split(newContent, EOL)

	if lines[len(lines)-1] == "" {
		lines = lines[:len(lines)-1]
	}

	if len(lines) > 0 && lines[0] == SnippetSeparator {
		lines = lines[1:]
	}

	newContent = strings.Join(lines, EOL)

	return ioutil.WriteFile(filePath, []byte(newContent), FileMode)
}

// Get all saved snippets
func GetAll(filePath string) ([]*Snippet, error) {
	return searchByKeyword("", filePath, SEARCH_MATCH_ANY)
}

// SearchExact exact search by keyword
func SearchExact(keyword string, filePath string) (*Snippet, error) {
	snippets, err := searchByKeyword(keyword, filePath, SEARCH_EXACT)

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
	return searchByKeyword(keyword, filePath, SEARCH_FUZZY)
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

	dir := path.Join(currentUser.HomeDir, ".snip")

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

	panicIfError(err)

	if _, err = os.Stat(dir); os.IsNotExist(err) {
		err = os.MkdirAll(dir, 0755)
		panicIfError(err)
	}

	filePath, err := SnippetFile()

	panicIfError(err)

	_, err = os.Stat(filePath)

	if os.IsNotExist(err) {
		file, err := os.OpenFile(filePath, os.O_RDONLY|os.O_CREATE, FileMode)
		panicIfError(err)
		defer util.Check(file.Close)
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

func searchByKeyword(keyword string, filePath string, exact SearchCode) ([]*Snippet, error) {
	var snippets []*Snippet

	file, err := os.Open(filePath)

	if err != nil {
		return snippets, err
	}

	scanner := getScanner(file)

	matcher := fuzzyMatcher

	if exact == SEARCH_EXACT {
		matcher = exactMatcher
	} else if exact == SEARCH_MATCH_ANY {
		matcher = func(k, c string) bool { return true }
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

		if exact == SEARCH_EXACT {
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

func panicIfError(e error) {
	if e != nil {
		panic(e)
	}
}
