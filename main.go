package main

import (
	"fmt"
	"os"
	"regexp"
	"strings"
)

func main() {
	text, _ := os.ReadFile("testdata/tolstoy/6157-0.txt")
	stext := string(text)

	fp, err := os.Create("/tmp/tolstoy.txt")
	defer fp.Close()
	if err != nil {
		fmt.Println(err)
	}

	
	re := regexp.MustCompile(`\b[\w']+\b`)
	for _, word := range re.FindAllString(stext, -1) {
		
		squashme := []string{"'s", "'t"}
		for _, suffix := range squashme {
			if strings.HasSuffix(word, suffix) {
				word = word[:len(word)-len(suffix)]
				word = word + suffix[1:]
			}
		}
		
		deleteme := []string{"'ll", "'m"}
		for _, suffix := range deleteme {
			if strings.HasSuffix(word, suffix) {
				word, _ = strings.CutSuffix(word, suffix)
			}
		}

		p := strings.Index(word, "'")
		if p != -1 {
			word = word[:p]
		}

		fmt.Fprintln(fp, word)
	}
}
