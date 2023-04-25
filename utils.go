package magictext

import (
	"crypto/md5"
	"fmt"
	"strconv"
	"time"
)

func CountTokens(text string) int {
	tokens := TikToken.Encode(text, nil, nil)
	return len(tokens)
}

func validateTokens(text string, maximum int) (int, bool) {
	numOfTokens := CountTokens(text)
	if numOfTokens > maximum {
		return numOfTokens, false
	}

	return numOfTokens, true
}

func randomFilename() string {
	return strconv.FormatInt(time.Now().UnixNano(), 10)
}

func hashString(text string) string {
	return fmt.Sprintf("%x", md5.Sum([]byte(text)))
}
