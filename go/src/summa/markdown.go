package summa

import (
	"bytes"
	"markdown"
	"os"
)

// markdownParse will take a markdown formatted input string and return
// HTML generated from parsing the markdown syntax
func markdownParse(input string) string {
	var outputBuf bytes.Buffer

	inputBuf := bytes.NewBufferString(input)
	p := markdown.NewParser(nil)
	p.Markdown(inputBuf, markdown.ToHTML(&outputBuf))

	return outputBuf.String()
}

// markdownParseFile will take the path to a markdown formatted file and
// return HTML generated from parsing the markdown syntax
func markdownParseFile(filePath string) (string, error) {
	var outputBuf bytes.Buffer

	f, err := os.Open(filePath)
	if err != nil {
		return "", err
	}

	p := markdown.NewParser(nil)
	p.Markdown(f, markdown.ToHTML(&outputBuf))

	return outputBuf.String(), nil
}
