package magictext

import (
	"fmt"

	"github.com/pkoukk/tiktoken-go"
)

var ttk *tiktoken.Tiktoken

func MakeFilename(infile, ext string) string {
	return fmt.Sprintf("%s.%s", infile, ext)
}

func CountTokens(text string) int {
	if ttk == nil {
		encoding := "cl100k_base" // support models: gpt-4, gpt-3.5-turbo, text-embedding-ada-002
		tke, err := tiktoken.GetEncoding(encoding)
		if err != nil {
			panic(err)
		}
		ttk = tke
	}

	token := ttk.Encode(text, nil, nil)
	return len(token)
}
