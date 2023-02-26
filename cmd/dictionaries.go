package main

import (
	"flag"
	"fmt"
	"os"
	"sort"

	dict "github.com/philhanna/dictionaries"
)

func main() {

	flag.Usage = func() {
		usage := `usage: dictionaries [OPTIONS] <inputFile> <outputFile>

Creates a list of words and frequency counts from an input text file.

positional parameters:
  inputFile        the input file
  outputFile       the output file

options are:
  -h, --help       Show this help text and exit
  -d, --debug      Output words and counts for manual inspection
`
		fmt.Fprintf(os.Stderr, "%s\n", usage)
	}
	var debug bool
	flag.BoolVar(&debug, "debug", false, "Output words and counts for manual inspection")
	flag.BoolVar(&debug, "d", false, "(short form of --debug)")
	flag.Parse()

	if flag.NArg() < 1 {
		fmt.Fprintf(os.Stderr, "No input file specified\n")
		return
	}
	inputFile := flag.Arg(0)

	if flag.NArg() < 2 {
		fmt.Fprintf(os.Stderr, "No output file specified\n")
		return
	}
	outputFile := flag.Arg(1)

	// Read the whole file as a string
	text, err := os.ReadFile(inputFile)
	if err != nil {
		fmt.Printf("Could not open input file %q: %s\n", inputFile, err.Error())
		return
	}

	// Parse the data
	ch := make(chan *dict.WordAndCount)
	go dict.ParseText(string(text), ch)
	defer close(ch)

	// Open the output file
	fpout, err := os.Create(outputFile)
	defer fpout.Close()
	if err != nil {
		fmt.Println(err)
		return
	}

	// Read from the parsing channel and write to the output file
	wordCount := 0
	switch debug {
	case true:
		for wac := range ch {
			if wac == nil {
				break
			}
			wordCount++
			word := wac.Word
			count := wac.Count
			fmt.Fprintf(fpout, "%s,%d\n", word, count)
		}
	default:
		words := make([]string, 0)
		for wac := range ch {
			if wac == nil {
				break
			}
			words = append(words, wac.Word)
		}
		sort.Slice(words, func(i, j int) bool {
			wordi := words[i]
			wordj := words[j]
			if len(wordi) > len(wordj) {
				return true
			}
			if len(wordi) < len(wordj) {
				return false
			}
			if wordi < wordj {
				return true
			}
			if wordi > wordj {
				return false
			}
			return false
		})
		for _, word := range words {
			wordCount++
			fmt.Fprintf(fpout, "%s,%d\n", word, len(word))
		}

	}
	fmt.Printf("%d words written to %s\n", wordCount, outputFile)
}
