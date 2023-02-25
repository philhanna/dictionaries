package main

import (
	"fmt"
	"os"
	"regexp"
	"sort"
	"strings"
)

type WordAndCount struct {
	Word  string
	Count int
}

func ParseText(text string, ch chan *WordAndCount) {

	// Create a map of words to their counts
	wordMap := make(map[string]int)

	// Make lists of apostrophe words and what to do about them.
	squashMe := []string{"'s", "'t"}
	deleteMe := []string{"'ll", "'m"}
	cutMe := []string{"'"}

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

		// Only words of length >= 3
		if len(word) < 3 {
			continue
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

	// Make a list of the keys and sort it in order of count descending.
	keys := []string{}
	for k := range wordMap {
		keys = append(keys, k)
	}
	sort.Slice(keys, func(i, j int) bool {
		iCount := wordMap[keys[i]]
		jCount := wordMap[keys[j]]
		return iCount > jCount
	})

	for _, word := range keys {
		count := wordMap[word]
		wac := WordAndCount{word, count}
		ch <- &wac
	}

	ch <- nil
}

func main() {

	// Read the whole file as a string
	text, _ := os.ReadFile("testdata/shakespeare/pg100.txt")

	// Parse the data
	ch := make(chan *WordAndCount)
	go ParseText(string(text), ch)
	defer close(ch)

	// Open the output file
	fpout, err := os.Create("/tmp/words.txt")
	defer fpout.Close()
	if err != nil {
		fmt.Println(err)
		return
	}

	// Read from the parsing channel and write to the output file

	for wac := range ch {
		if wac == nil {
			break
		}
		word := wac.Word
		count := wac.Count
		fmt.Fprintf(fpout, "%s,%d\n", word, count)
	}
}
