package main

import (
	"fmt"
	"os"
	"regexp"
	"sort"
	"strings"
)

func main() {

	// Create a map of words to their counts
	wordMap := make(map[string]int)

	// Make lists of apostrophe words and what to do about them.
	squashMe := []string{"'s", "'t"}
	deleteMe := []string{"'ll", "'m"}
	cutMe := []string{"'"}
	
	// Form the regular expression
	re := regexp.MustCompile(`\b[\w']+\b`)

	// Read the whole file as a string
	text, _ := os.ReadFile("testdata/tolstoy/6157-0.txt")
	stext := string(text)

	// Read the text and extract the words
	for _, word := range re.FindAllString(stext, -1) {
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

	// Open the output file
	fpout, err := os.Create("/tmp/tolstoy.txt")
	defer fpout.Close()
	if err != nil {
		fmt.Println(err)
		return
	}

	// Write the map
	for _, word := range keys {
		count := wordMap[word]
		fmt.Fprintf(fpout, "%s,%d\n", word, count)
	}

}
