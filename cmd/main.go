package main

import (
	"flag"
	"fmt"
	"os"

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
  -h               Show this help text and exit
`
		fmt.Fprintf(os.Stderr, "%s\n", usage)
	}
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
	for wac := range ch {
		if wac == nil {
			break
		}
		word := wac.Word
		count := wac.Count
		fmt.Fprintf(fpout, "%s,%d\n", word, count)
	}
}
