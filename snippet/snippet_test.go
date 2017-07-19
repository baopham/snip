package snippet_test

import (
	"fmt"
	. "gopkg.in/baopham/snip.v2/snippet"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"io/ioutil"
	"os"
	"path"
	"strings"
)

type SnippetContent string

func (c SnippetContent) mustContain(snippet Snippet) {
	content := string(c)
	Expect(content).To(ContainSubstring(snippet.Keyword))
	Expect(content).To(ContainSubstring(snippet.Content))
	Expect(content).To(ContainSubstring(snippet.Description))
}

func (c SnippetContent) mustNotContain(snippet Snippet) {
	content := string(c)
	Expect(content).ToNot(ContainSubstring(snippet.Keyword))
	Expect(content).ToNot(ContainSubstring(snippet.Content))
	Expect(content).ToNot(ContainSubstring(snippet.Description))
}

var _ = Describe("Snippet", func() {
	var (
		fakeFilePath   string
		snippetCounter = 1
	)

	saveSnippet := func(snippet Snippet, filePath string) {
		err := snippet.Save(filePath)
		Expect(err).To(BeNil())
	}

	removeFakeFile := func(fakeFilePath string) {
		if _, err := os.Stat(fakeFilePath); err == nil {
			err = os.Remove(fakeFilePath)
			if err != nil {
				panic(err)
			}
		}
	}

	getFileContent := func(fakeFilePath string) SnippetContent {
		b, err := ioutil.ReadFile(fakeFilePath)
		Expect(err).To(BeNil())
		str := strings.TrimSpace(string(b))
		return SnippetContent(str)
	}

	seedSnippet := func() Snippet {
		counter := fmt.Sprint(snippetCounter)
		snippet := Snippet{
			Keyword:     "keyword " + counter,
			Description: "description " + counter,
			Content:     "content " + counter,
		}
		snippetCounter++
		return snippet
	}

	BeforeEach(func() {
		dir, err := os.Getwd()

		if err != nil {
			panic(err)
		}

		fakeFilePath = path.Join(dir, ".snippets")

		removeFakeFile(fakeFilePath)
	})

	AfterEach(func() {
		removeFakeFile(fakeFilePath)
	})

	Context("when calling snippet.Save()", func() {
		It("should save the snippet", func() {
			By("appending the snippet to file ~/.snip/snippets")

			snippet := seedSnippet()

			saveSnippet(snippet, fakeFilePath)

			content := getFileContent(fakeFilePath)

			content.mustContain(snippet)
		})
	})

	Context("when saving multiple snippets", func() {
		It("should save the snippets", func() {
			snippet1 := seedSnippet()

			saveSnippet(snippet1, fakeFilePath)

			snippet2 := seedSnippet()

			saveSnippet(snippet2, fakeFilePath)

			content := strings.Split(string(getFileContent(fakeFilePath)), "\n")

			SnippetContent(content[0]).mustContain(snippet1)
			SnippetContent(content[1]).mustContain(snippet2)
		})
	})

	Context("when snippet having multiple lines", func() {
		It("should save the snippet correctly", func() {
			snippet := seedSnippet()

			saveSnippet(snippet, fakeFilePath)

			content := getFileContent(fakeFilePath)

			content.mustContain(snippet)
		})
	})

	Context("when calling SearchExact() by an existing keyword", func() {
		It("should return the found Snippet", func() {
			snippet := seedSnippet()

			saveSnippet(snippet, fakeFilePath)

			found, err := SearchExact(snippet.Keyword, fakeFilePath)

			Expect(err).To(BeNil())
			Expect(*found).To(Equal(snippet))
		})
	})

	Context("when calling Search()", func() {
		var snippet1, snippet2, snippet3 Snippet

		BeforeEach(func() {
			snippet1 = Snippet{
				Keyword:     "port",
				Description: "Find processes using a certain port",
				Content:     "lsof -i :{p}",
			}
			saveSnippet(snippet1, fakeFilePath)
			snippet2 = Snippet{
				Keyword:     "port2",
				Description: "Find processes using a certain port",
				Content:     "lsof -i :{p}",
			}
			saveSnippet(snippet2, fakeFilePath)
			snippet3 = seedSnippet()
			saveSnippet(snippet3, fakeFilePath)
		})

		assertSearchResult := func(searchTerm string) {
			found, err := Search(searchTerm, fakeFilePath)
			Expect(err).To(BeNil())

			Expect(found).To(HaveLen(2))
			Expect(*found[0]).To(Equal(snippet1))
			Expect(*found[1]).To(Equal(snippet2))
		}

		It("should return all the found snippets that fuzzy match the search term", func() {
			By("searching by substring of keyword")

			assertSearchResult("port")

			By("searching by substring of description")

			assertSearchResult("find process")
		})
	})

	Context("when saving the same snippet (same keyword)", func() {
		It("should not save it again", func() {
			snippet := seedSnippet()

			saveSnippet(snippet, fakeFilePath)

			err := snippet.Save(fakeFilePath)

			By("returning SnippetAlreadyExistError error")

			Expect(err).To(MatchError(SnippetAlreadyExistError{Keyword: snippet.Keyword}))
		})
	})

	Context("when calling snippet.Build()", func() {
		It("should build the snippet using the given placeholders", func() {
			By("replacing one placeholder")

			snippet := Snippet{
				Keyword:     "port",
				Description: "Find processes using a certain port",
				Content:     "lsof -i :{p}",
			}

			content := snippet.Build(map[string]string{
				"p": "9001",
			})

			Expect(content).To(Equal("lsof -i :9001"))

			By("replacing multiple placeholders including duplicates")

			snippet = Snippet{
				Keyword:     "foo",
				Description: "Foo bar",
				Content:     "{a} foo {b} bar {c} {a} {b} {c}",
			}

			content = snippet.Build(map[string]string{
				"a": "A",
				"b": "B",
				"c": "C",
			})

			Expect(content).To(Equal("A foo B bar C A B C"))
		})
	})

	Context("when calling snippet.Remove()", func() {
		saveThreeSnippets := func() (Snippet, Snippet, Snippet) {
			By("saving 3 snippets")

			snippet1 := seedSnippet()

			saveSnippet(snippet1, fakeFilePath)

			snippet2 := seedSnippet()

			saveSnippet(snippet2, fakeFilePath)

			snippet3 := seedSnippet()

			saveSnippet(snippet3, fakeFilePath)

			content := getFileContent(fakeFilePath)

			By("checking that the snippets are indeed saved")

			content.mustContain(snippet1)
			content.mustContain(snippet2)
			content.mustContain(snippet3)

			return snippet1, snippet2, snippet3
		}

		Context("when removing the first saved snippet", func() {
			It("should remove the snippet from the file", func() {
				snippet1, snippet2, snippet3 := saveThreeSnippets()

				By("removing the first snippet that was saved")

				err := snippet1.Remove(fakeFilePath)

				Expect(err).To(BeNil())

				By("checking that now the first snippet is removed")

				newContent := getFileContent(fakeFilePath)

				Expect(strings.Split(string(newContent), "\n")).To(HaveLen(2))
				newContent.mustNotContain(snippet1)
				newContent.mustContain(snippet2)
				newContent.mustContain(snippet3)
			})
		})

		Context("when removing the last saved snippet", func() {
			It("should remove the snippet from the file", func() {
				snippet1, snippet2, snippet3 := saveThreeSnippets()

				By("removing the last snippet that was saved")

				err := snippet3.Remove(fakeFilePath)

				Expect(err).To(BeNil())

				By("checking that now the last snippet is removed")

				newContent := getFileContent(fakeFilePath)

				Expect(strings.Split(string(newContent), "\n")).To(HaveLen(2))
				newContent.mustContain(snippet1)
				newContent.mustContain(snippet2)
				newContent.mustNotContain(snippet3)
			})
		})

		Context("when removing the middle saved snippet", func() {
			It("should remove the snippet from the file", func() {
				snippet1, snippet2, snippet3 := saveThreeSnippets()

				By("removing the middle snippet that was saved")

				err := snippet2.Remove(fakeFilePath)

				Expect(err).To(BeNil())

				By("checking that now the middle snippet is removed")

				newContent := getFileContent(fakeFilePath)

				newContent.mustContain(snippet1)
				newContent.mustNotContain(snippet2)
				newContent.mustContain(snippet3)
			})
		})
	})
})
