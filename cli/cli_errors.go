package cli

type MissingArgumentError struct {
	Message string
}

func (e MissingArgumentError) Error() string {
	return e.Message
}
