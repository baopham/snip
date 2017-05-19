package snippet_test

import (
	"fmt"
	. "github.com/baopham/snip/snippet"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"io/ioutil"
	"os"
	"path"
	"strings"
)

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

	getFileContent := func(fakeFilePath string) string {
		b, err := ioutil.ReadFile(fakeFilePath)
		Expect(err).To(BeNil())
		return string(b)
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

			Expect(content).To(ContainSubstring(snippet.String()))
		})
	})

	Context("when saving multiple snippets", func() {
		It("should save the snippets with correct separator", func() {
			snippet1 := seedSnippet()

			saveSnippet(snippet1, fakeFilePath)

			snippet2 := seedSnippet()

			saveSnippet(snippet2, fakeFilePath)

			content := strings.Split(getFileContent(fakeFilePath), EOL+SnippetSeparator+EOL)

			Expect(content[0]).To(Equal(snippet1.String()))
			Expect(content[1]).To(Equal(snippet2.String()))
		})
	})

	Context("when snippet having multiple lines", func() {
		It("should save the snippet correctly", func() {
			snippet := seedSnippet()

			saveSnippet(snippet, fakeFilePath)

			content := getFileContent(fakeFilePath)

			Expect(content).To(ContainSubstring(snippet.String()))
			Expect(content).To(ContainSubstring(snippet.Content))
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

	Context("when calling Search() by a keyword", func() {
		It("should return all the found snippets that fuzzy match the keyword", func() {
			snippet1 := Snippet{
				Keyword:     "port",
				Description: "Find processes using a certain port",
				Content:     "lsof -i :{p}",
			}

			saveSnippet(snippet1, fakeFilePath)

			snippet2 := Snippet{
				Keyword:     "port2",
				Description: "Find processes using a certain port",
				Content:     "lsof -i :{p}",
			}

			saveSnippet(snippet2, fakeFilePath)

			snippet3 := seedSnippet()

			saveSnippet(snippet3, fakeFilePath)

			found, err := Search("port", fakeFilePath)
			Expect(err).To(BeNil())

			Expect(found).To(HaveLen(2))
			Expect(*found[0]).To(Equal(snippet1))
			Expect(*found[1]).To(Equal(snippet2))
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
		separator := EOL + SnippetSeparator + EOL

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

			Expect(content).To(Equal(snippet1.String() + separator + snippet2.String() + separator + snippet3.String()))

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

				Expect(newContent).To(Equal(snippet2.String() + separator + snippet3.String()))
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

				Expect(newContent).To(Equal(snippet1.String() + separator + snippet2.String()))
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

				Expect(newContent).To(Equal(snippet1.String() + separator + snippet3.String()))
			})
		})
	})
})
