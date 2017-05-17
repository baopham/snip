package snippet_test

import (
	. "github.com/baopham/snippets-cli/snippet"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"io/ioutil"
	"os"
	"path"
	"strings"
)

var _ = Describe("Snippet", func() {
	var (
		fakeFilePath string
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
			By("appending the snippet to file ~/.snippets-cli/snippets")

			snippet := Snippet{
				Keyword:     "port",
				Description: "Find processes using a certain port",
				Content:     "lsof -i :{p}",
			}

			saveSnippet(snippet, fakeFilePath)

			b, err := ioutil.ReadFile(fakeFilePath)
			Expect(err).To(BeNil())

			Expect(string(b)).To(ContainSubstring(snippet.String()))
		})
	})

	Context("when saving multiple snippets", func() {
		It("should save the snippets with correct separator", func() {
			snippet1 := Snippet{
				Keyword:     "port",
				Description: "Find processes using a certain port",
				Content:     "lsof -i :{p}",
			}

			saveSnippet(snippet1, fakeFilePath)

			snippet2 := Snippet{
				Keyword:     "docker-ssh",
				Description: "SSH into docker container",
				Content:     "docker exec -it {id} bash",
			}

			saveSnippet(snippet2, fakeFilePath)

			b, err := ioutil.ReadFile(fakeFilePath)
			Expect(err).To(BeNil())

			content := strings.Split(string(b), EOL+SnippetSeparator+EOL)

			Expect(content[0]).To(Equal(snippet1.String()))
			Expect(content[1]).To(Equal(snippet2.String()))
		})
	})

	Context("when snippet having multiple lines", func() {
		It("should save the snippet correctly", func() {
			snippet := Snippet{
				Keyword:     "foo",
				Description: "Some long command",
				Content: `foo
				bar
				foobar
				`,
			}

			saveSnippet(snippet, fakeFilePath)

			b, err := ioutil.ReadFile(fakeFilePath)
			Expect(err).To(BeNil())

			content := string(b)

			Expect(content).To(ContainSubstring(snippet.String()))
			Expect(content).To(ContainSubstring(snippet.Content))
		})
	})

	Context("when calling SearchExact() by an existing keyword", func() {
		It("should return the found Snippet", func() {
			snippet := Snippet{
				Keyword:     "docker-ssh",
				Description: "SSH into docker container",
				Content:     "docker exec -it {id} bash",
			}

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
				Keyword:     "port-v2",
				Description: "Find processes using a certain port",
				Content:     "lsof -i :{p}",
			}

			saveSnippet(snippet2, fakeFilePath)

			snippet3 := Snippet{
				Keyword:     "foo",
				Description: "Foo bar",
				Content:     "foo bar",
			}

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
			snippet := Snippet{
				Keyword:     "port",
				Description: "Find processes using a certain port",
				Content:     "lsof -i :{p}",
			}

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
})
