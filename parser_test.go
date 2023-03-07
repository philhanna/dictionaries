package dictionaries

import (
	"fmt"
	"os"
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
			for wac := range ParseText(tt.input) {
				have = append(have, wac.Word)
			}
			assert.Equal(t, tt.want, have)
		})
	}
}

func TestParseWebPage(t *testing.T) {
	t.Skip()
	url := "https://www.cnn.com"
	fp, _ := os.Create("a.txt")
	defer fp.Close()
	for wac := range ParseWebPage(url) {
		fmt.Fprintln(fp, wac)
	}
}
