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

func randString(min, max int) string {
	rand.Seed(time.Now().UnixNano())

	var chineseChars = []rune("的一是在不了有和人这中大为上个国我以要他时来用们生到作地于出就分对成会可主发年动同工也能下过子说产种面而方后多定行学法所民得经十三之进着等部度家电力力里如水化高自二理起小物现实加量都两体制机当使点从业本去把性好应开它合还因由其些然前外天政四日那社义事平形相全表间样与关各重新线内数正心反你明看原又么利比或但质气第向道命此变条只没结解问意建月公无系军很情者最立代想已通并提直题党程展五果料象员革位入常文总次品式活设及管特件长求老头基资较新青岛先安先河各式样石紫军新村明园广场等地")
	b := make([]rune, rand.Intn(max-min)+min)
	for i := range b {
		b[i] = chineseChars[rand.Intn(len(chineseChars))]
	}
	return string(b)
}

func hashString(text string) string {
	return fmt.Sprintf("%x", md5.Sum([]byte(text)))
}
