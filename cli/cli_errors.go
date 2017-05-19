package cli

type MissingInfoError struct {
	Message string
}

func (e MissingInfoError) Error() string {
	return e.Message
}

