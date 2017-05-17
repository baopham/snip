package snippets

import (
	"fmt"
)

// SnippetAlreadyExistError error when the snippet already exists
type SnippetAlreadyExistError struct {
	Keyword string
}

func (e SnippetAlreadyExistError) Error() string {
	return fmt.Sprintf("Snippet %s already exists", e.Keyword)
}
