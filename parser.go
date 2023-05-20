package dictionaries

import (
	"io"
	"log"
	"net/http"
	"regexp"
	"sort"
	"strings"

	"golang.org/x/net/html"
)

// ---------------------------------------------------------------------
// Type definitions
// ---------------------------------------------------------------------

// WordAndCount keeps a word string and the number of occurrences it has
// in the source text.
type WordAndCount struct {
	Word  string
	Count int
}

// ---------------------------------------------------------------------
// Functions
// ---------------------------------------------------------------------

// ParseText is a Python-like generator that reads through the text,
// parsing individual words, and keeping track of them in a map of word
// to number of occurrences. Then it sorts the map in descending order
// of the word count and yields its items back to the caller.
func ParseText(text string) <-chan string {

	// Create a map of words to their counts
	wordMap := make(map[string]int)

	// Make lists of apostrophe words and what to do about them.
	squashMe := []string{"'s", "'t"}  // Remove the apostrophe
	deleteMe := []string{"'ll", "'m"} // Chop the suffix
	cutMe := []string{"'"}            // Remove everything from apostrophe to end

	// Form the regular expression
	re := regexp.MustCompile(`\b[a-zA-Z']+\b`)

	// Read the text and extract the words
	for _, word := range re.FindAllString(text, -1) {
		for _, suffix := range squashMe {
			if strings.HasSuffix(word, suffix) {
				word = word[:len(word)-len(suffix)]
				word = word + suffix[1:]
			}
		}
		for _, suffix := range deleteMe {
			if strings.HasSuffix(word, suffix) {
				word, _ = strings.CutSuffix(word, suffix)
			}
		}
		for _, suffix := range cutMe {
			p := strings.Index(word, suffix)
			if p != -1 {
				word = word[:p]
			}
		}

		// Convert to uppercase
		word = strings.ToUpper(word)

		// Add to map
		_, ok := wordMap[word]
		if ok {
			wordMap[word]++
		} else {
			wordMap[word] = 1
		}
	}

	// Make a list of the keys and sort it alphabetically
	keys := []string{}
	for k := range wordMap {
		keys = append(keys, k)
	}
	sort.Slice(keys, func(i, j int) bool {
		return keys[i] < keys[j]
	})

	// Send the resulting sorted list up the channel
	ch := make(chan string)
	go func() {
		defer close(ch)
		for _, word := range keys {
			ch <- word
		}
	}()
	return ch
}

// ParseWebPage is a Python-like generator that downloads the contents of a
// web page and yields its words.
func ParseWebPage(url string) <-chan string {

	// Get the text from the web page, logging any errors
	resp, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	// Get the actual text
	b, _ := io.ReadAll(resp.Body)
	body := string(b)
	node, _ := html.Parse(strings.NewReader(body))
	text := getText(node)

	return ParseText(text)
}

func getText(node *html.Node) string {
	sb := ""
	if node.Type == html.TextNode {
		sb += " "
		sb += node.Data
	} else {
		tagName := node.Data
		if tagName == "head" || tagName == "script" || tagName == "style" {
			// Skip this
		} else {
			for child := node.FirstChild; child != nil; child = child.NextSibling {
				sb += " "
				sb += getText(child)
			}
		}
	}
	return sb
}
