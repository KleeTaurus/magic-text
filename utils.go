package magictext

import (
	"crypto/md5"
	"fmt"
	"math/rand"
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

func randFilename() string {
	return strconv.FormatInt(time.Now().UnixNano(), 10)
}

func randString() string {
	rand.Seed(time.Now().UnixNano())

	var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	b := make([]rune, rand.Intn(60)+80)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func hashString(text string) string {
	return fmt.Sprintf("%x", md5.Sum([]byte(text)))
}
