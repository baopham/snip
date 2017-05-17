package snippets_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestSnippets(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Snippets Suite")
}
