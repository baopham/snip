package snippet

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"os/user"
	"path"
	"strings"

	"github.com/baopham/snip/util"
	"github.com/renstrom/fuzzysearch/fuzzy"
)

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

	defer util.Check(file.Close)

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

	w := csv.NewWriter(file)

	if err := w.Write([]string{s.Keyword, s.Content, s.Description}); err != nil {
		return err
	}

	w.Flush()

	return w.Error()
}

// Remove a saved snippet
func (s *Snippet) Remove(filePath string) error {
	rows := make([][]string, 0)

	file, err := os.Open(filePath)

	if err != nil {
		return err
	}

	csvr := csv.NewReader(file)

	for {
		row, err := csvr.Read()

		if err == io.EOF {
			break
		}

		if err != nil {
			return err
		}

		keyword, content, description := row[0], row[1], row[2]

		if s.Keyword == keyword {
			continue
		}

		rows = append(rows, []string{keyword, content, description})
	}

	err = file.Close()

	if err != nil {
		return err
	}

	file, err = os.OpenFile(filePath, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, FileMode)

	defer util.Check(file.Close)

	w := csv.NewWriter(file)

	return w.WriteAll(rows)
}

// Get all saved snippets
func GetAll(filePath string) ([]*Snippet, error) {
	return searchSnippets("", filePath, SEARCH_MATCH_ANY)
}

// SearchExact exact search by keyword
func SearchExact(keyword, filePath string) (*Snippet, error) {
	snippets, err := searchSnippets(keyword, filePath, SEARCH_EXACT)

	if err != nil {
		return nil, err
	}

	if len(snippets) == 0 {
		return nil, nil
	}

	return snippets[0], nil
}

// Search fuzzy search by given search term
func Search(searchTerm, filePath string) ([]*Snippet, error) {
	return searchSnippets(searchTerm, filePath, SEARCH_FUZZY)
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

	filePath := path.Join(dir, "snippets.csv")

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

func searchSnippets(searchTerm, filePath string, exact SearchCode) ([]*Snippet, error) {
	var snippets []*Snippet

	file, err := os.Open(filePath)
	defer util.Check(file.Close)

	if err != nil {
		return snippets, err
	}

	matcher := fuzzy.MatchFold

	if exact == SEARCH_EXACT {
		matcher = exactMatcher
	} else if exact == SEARCH_MATCH_ANY {
		matcher = func(k, c string) bool { return true }
	}

	csvr := csv.NewReader(file)

	for {
		row, err := csvr.Read()

		if err == io.EOF {
			return snippets, nil
		}

		if err != nil {
			return snippets, err
		}

		keyword, content, description := row[0], row[1], row[2]

		if !matcher(searchTerm, keyword) && !matcher(searchTerm, content) && !matcher(searchTerm, description) {
			continue
		}

		found := &Snippet{
			Keyword:     keyword,
			Content:     content,
			Description: description,
		}

		snippets = append(snippets, found)

		if exact == SEARCH_EXACT {
			return snippets, nil
		}
	}

	return snippets, nil
}

func exactMatcher(source string, target string) bool {
	return strings.TrimSpace(source) == strings.TrimSpace(target)
}

func panicIfError(e error) {
	if e != nil {
		panic(e)
	}
}
