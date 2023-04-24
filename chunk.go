package magictext

import (
	"crypto/md5"
	"fmt"
	"strings"
	"time"
	"unicode/utf8"
)

type Category int8

const (
	CatText Category = iota + 1
	CatSummary
)

type Chunk struct {
	Category Category `json:"category"`
	Depth    int      `json:"depth"`
	Seq      int      `json:"seq"`
	ParentID string   `json:"parent_id"`
	ID       string   `json:"id"`
	Text     string   `json:"text"`
	Tokens   int      `json:"tokens"`
}

type ChunkSlice []*Chunk

func NewTextChunk(seq int, text string) *Chunk {
	return &Chunk{
		Category: CatText,
		Depth:    0,
		Seq:      seq,
		Text:     text,
		ID:       fmt.Sprintf("%x", md5.Sum([]byte(text))),
		Tokens:   CountTokens(text),
	}
}

func NewSummaryChunk(seq int, summary string, depth int) *Chunk {
	return &Chunk{
		Category: CatSummary,
		Depth:    depth,
		Seq:      seq,
		Text:     summary,
		ID:       fmt.Sprintf("%x", md5.Sum([]byte(summary))),
		Tokens:   CountTokens(summary),
	}
}

func (c *Chunk) String() string {
	return fmt.Sprintf("%d:%d:%d:%s:%s:%d", c.Category, c.Depth, c.Seq, c.ParentID, c.ID, c.Tokens)
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

func (cs ChunkSlice) SubGroups(maxTokensPerRequest, maxChunksInGroup int) []ChunkSlice {
	groups := make([]ChunkSlice, 0, len(cs)/2)
	chunkGroup := make(ChunkSlice, 0, maxChunksInGroup)
	for _, chunk := range cs {
		if chunkGroup.Tokens()+chunk.Tokens > maxTokensPerRequest || len(chunkGroup) >= maxChunksInGroup {
			groups = append(groups, chunkGroup)
			// reset chunkGroup to empty
			chunkGroup = make(ChunkSlice, 0, maxChunksInGroup)
		}
		chunkGroup = append(chunkGroup, chunk)
	}
	groups = append(groups, chunkGroup)
	return groups
}

func (cs ChunkSlice) String() string {
	cm := make(map[Category]bool)
	dm := make(map[int]bool)

	for _, chunk := range cs {
		if _, ok := cm[chunk.Category]; !ok {
			cm[chunk.Category] = true
		}

		if _, ok := dm[chunk.Depth]; !ok {
			dm[chunk.Depth] = true
		}
	}
	return fmt.Sprintf("Category: %v, Depth: %v, Childs: %d, Total Tokens: %d", cm, dm, len(cs), cs.Tokens())
}

func NewTextChunk2(seq int, text string) TextChunk {
	text = strings.TrimSpace(text)

	tc := TextChunk{}
	tc.Seq = seq
	tc.Text = text
	tc.ID = fmt.Sprintf("%x", md5.Sum([]byte(text)))
	tc.Tokens = CountTokens(text)

	return tc
}

type TextChunk struct {
	ID     string `json:"id"`
	Seq    int    `json:"seq"`
	Text   string `json:"text"`
	Tokens int    `json:"tokens"`
}

func NewCaptionChunk(seq int, text string, from time.Time) CaptionChunk {
	text = strings.TrimSpace(text)

	cc := CaptionChunk{}
	cc.Seq = seq
	cc.Text = text
	cc.From = from
	cc.ID = fmt.Sprintf("%x", md5.Sum([]byte(text)))
	cc.Tokens = CountTokens(text)

	return cc
}

type CaptionChunk struct {
	From time.Time `json:"from"`
	TextChunk
}

func (c CaptionChunk) String() string {
	text := c.Text
	maxLength := 80
	if utf8.RuneCountInString(c.Text) > maxLength {
		text = fmt.Sprintf("%s...", string([]rune(c.Text)[:maxLength-3]))
	}

	return fmt.Sprintf("%s <%04d> %s %s", c.ID[:8], c.Tokens, c.From.Format("15:04:05"), text)
}
