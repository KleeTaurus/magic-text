package textsummary

import (
	"crypto/md5"
	"fmt"
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

func (cd *Chunk) Len() int {
	return len(cd.Text)
}

func (cd *Chunk) RuneCountInString() int {
	return utf8.RuneCountInString(cd.Text)
}

func (cd *Chunk) String() string {
	return fmt.Sprintf("%d:%d:%s:%s:%d",
		cd.Category, cd.Depth, cd.ParentID, cd.ID, cd.RuneCountInString())
}
