package magictext

import (
	"fmt"
)

func MakeFilename(infile, ext string) string {
	return fmt.Sprintf("%s.%s", infile, ext)
}

func CountTokens(text string) int {
	token := TikToken.Encode(text, nil, nil)
	return len(token)
}

func validateTokens(text string, maximum int) (int, bool) {
	tokens := CountTokens(text)
	if tokens > maximum {
		return tokens, false
	}

	return tokens, true
}
