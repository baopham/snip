package util

// Check defer
// http://stackoverflow.com/questions/40397781/gometalinter-errcheck-returns-a-warning-on-deferring-a-func-which-returns-a-va
func Check(f func() error) {
	if err := f(); err != nil {
		panic(err)
	}
}
