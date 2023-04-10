package summaryit

import (
	"fmt"
	"path"
	"strings"

	"github.com/pkoukk/tiktoken-go"
)

func GenerateFilename(infile, ext string) string {
	dirname := path.Dir(infile)
	basename := path.Base(infile)

	return fmt.Sprintf("%s/%s%s", dirname, strings.Split(basename, ".")[0], ext)
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
