package magictext

import (
	"fmt"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"
)

type Chunk struct {
	ID       string     `json:"id"`
	Seq      int        `json:"seq"`
	Height   int        `json:"height"`
	Text     string     `json:"text"`
	Tokens   int        `json:"tokens"`
	Children ChunkSlice `json:"children"`
}

type ChunkSlice []*Chunk

func NewChunk(seq int, text string) *Chunk {
	return &Chunk{
		Seq:    seq,
		Text:   text,
		ID:     hashString(text),
		Tokens: CountTokens(text),
	}
}

func (c *Chunk) String() string {
	return fmt.Sprintf("%s:%d:%d:%d", c.ID, c.Height, c.Seq, c.Tokens)
}

func (cs ChunkSlice) Tokens() int {
	tokens := 0
	for _, chunk := range cs {
		tokens += chunk.Tokens
	}
	return tokens
}

func (cs ChunkSlice) TokenString() string {
	s := strings.Builder{}
	for _, chunk := range cs {
		s.WriteString(chunk.Text + " ")
	}
	return s.String()
}

func (cs ChunkSlice) String() string {
	dm := make(map[int]bool)

	for _, chunk := range cs {
		if _, ok := dm[chunk.Height]; !ok {
			dm[chunk.Height] = true
		}
	}

	sb := strings.Builder{}
	for k := range dm {
		sb.WriteString(strconv.Itoa(k))
	}
	return fmt.Sprintf("Height: %s, Children: %d, Total Tokens: %d", sb.String(), len(cs), cs.Tokens())
}

type TextChunk struct {
	ID     string `json:"id"`
	Seq    int    `json:"seq"`
	Text   string `json:"text"`
	Tokens int    `json:"tokens"`
}

type CaptionChunk struct {
	From time.Time `json:"from"`
	TextChunk
}

func NewCaptionChunk(seq int, text string, from time.Time) CaptionChunk {
	text = strings.TrimSpace(text)

	cc := CaptionChunk{}
	cc.Seq = seq
	cc.Text = text
	cc.From = from
	cc.ID = hashString(text)
	cc.Tokens = CountTokens(text)

	return cc
}

func (c CaptionChunk) String() string {
	text := c.Text
	maxLength := 80
	if utf8.RuneCountInString(c.Text) > maxLength {
		text = fmt.Sprintf("%s...", string([]rune(c.Text)[:maxLength-3]))
	}

	return fmt.Sprintf("%s <%04d> %s %s", c.ID[:8], c.Tokens, c.From.Format("15:04:05"), text)
}
