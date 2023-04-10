package summaryit

import (
	"crypto/md5"
	"fmt"
	"strings"
	"unicode/utf8"
)

type Category int8

const (
	CatText Category = iota + 1
	CatSummary
)

type Chunk struct {
	Category Category `json:"category"`
	Depth    uint     `json:"depth"`
	ParentID string   `json:"parent_id"`
	ID       string   `json:"id"`
	Text     string   `json:"text"`
}

type ChunkSlice []*Chunk

func NewTextChunk(text string) *Chunk {
	return &Chunk{
		Depth:    0,
		Category: CatText,
		Text:     text,
		ID:       fmt.Sprintf("%x", md5.Sum([]byte(text))),
	}
}

func NewSummaryChunk(summary string, depth uint) *Chunk {
	return &Chunk{
		Depth:    depth,
		Category: CatSummary,
		Text:     summary,
		ID:       fmt.Sprintf("%x", md5.Sum([]byte(summary))),
	}
}

func (c *Chunk) Len() int {
	return len(c.Text)
}

func (c *Chunk) Tokens() int {
	return utf8.RuneCountInString(c.Text)
}

func (c *Chunk) String() string {
	return fmt.Sprintf("%d:%d:%s:%s:%d", c.Category, c.Depth, c.ParentID, c.ID, c.Tokens())
}

func (cs ChunkSlice) Tokens() int {
	c := 0
	for _, chunk := range cs {
		c += chunk.Tokens()
	}
	return c
}

func (cs ChunkSlice) TokenString() string {
	s := strings.Builder{}
	for _, chunk := range cs {
		s.WriteString(chunk.Text + " ")
	}
	return s.String()
}

func (cs ChunkSlice) Group(maxTokensPerRequest, maxChunksInGroup int) []ChunkSlice {
	groups := make([]ChunkSlice, 0, len(cs)/2)
	chunkGroup := make(ChunkSlice, 0, maxChunksInGroup)
	for _, c := range cs {
		if c.Tokens()+chunkGroup.Tokens() > maxTokensPerRequest || len(chunkGroup) >= maxChunksInGroup {
			groups = append(groups, chunkGroup)
			chunkGroup = make(ChunkSlice, 0, maxChunksInGroup)
		}
		chunkGroup = append(chunkGroup, c)
	}
	groups = append(groups, chunkGroup)
	return groups
}

func (cs ChunkSlice) String() string {
	c := 0
	cm := make(map[Category]bool)
	dm := make(map[uint]bool)

	for _, chunk := range cs {
		c++
		if _, ok := cm[chunk.Category]; !ok {
			cm[chunk.Category] = true
		}

		if _, ok := dm[chunk.Depth]; !ok {
			dm[chunk.Depth] = true
		}
	}
	return fmt.Sprintf("Category: %v, Depth: %v, Childs: %d, Total Tokens: %d", cm, dm, c, cs.Tokens())
}
