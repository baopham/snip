package snippet_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestSnippet(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Snippet Suite")
}
