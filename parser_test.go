package dictionaries

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseText(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  []string
	}{
		{"short text", "Now is the time", []string{"IS", "NOW", "THE", "TIME"}},
		{"empty text", "", []string{}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			have := make([]string, 0)
			for word := range ParseText(tt.input) {
				have = append(have, word)
			}
			assert.Equal(t, tt.want, have)
		})
	}
}

func TestParseWebPage(t *testing.T) {
	url := "https://www.cnn.com"
	tmp := os.TempDir()
	filename := filepath.Join(tmp, "parse_web_page.txt")
	log.Printf("Writing test output to %s\n", filename)
	fp, err := os.Create(filename)
	assert.Nil(t, err)
	defer fp.Close()
	ch := ParseWebPage(url)
	for word := range ch {
		fmt.Fprintln(fp, word)
	}
}
