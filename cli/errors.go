package cli

type MissingInfoError struct {
	Message string
}

type NotFoundSnippet struct {
	Keyword string
}

func (e NotFoundSnippet) Error() string {
	return "Could not find a snippet with keyword: " + e.Keyword
}

func (e MissingInfoError) Error() string {
	return e.Message
}
