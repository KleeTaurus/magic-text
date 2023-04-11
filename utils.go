package summaryit

import (
	"fmt"

	"github.com/pkoukk/tiktoken-go"
)

func MakeFilename(infile, ext string) string {
	return fmt.Sprintf("%s.%s", infile, ext)
}

func CountTokens(text string) int {
	encoding := "cl100k_base" // support models: gpt-4, gpt-3.5-turbo, text-embedding-ada-002
	tke, err := tiktoken.GetEncoding(encoding)
	if err != nil {
		panic(err)
	}

	token := tke.Encode(text, nil, nil)
	return len(token)
}
