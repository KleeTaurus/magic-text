package magictext

import (
	"crypto/md5"
	"fmt"
	"strings"
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
	Tokens   int      `json:"tokens"`
}

type ChunkSlice []*Chunk

func NewTextChunk(text string) *Chunk {
	return &Chunk{
		Depth:    0,
		Category: CatText,
		Text:     text,
		ID:       fmt.Sprintf("%x", md5.Sum([]byte(text))),
		Tokens:   CountTokens(text),
	}
}

func NewSummaryChunk(summary string, depth uint) *Chunk {
	return &Chunk{
		Depth:    depth,
		Category: CatSummary,
		Text:     summary,
		ID:       fmt.Sprintf("%x", md5.Sum([]byte(summary))),
		Tokens:   CountTokens(summary),
	}
}

func (c *Chunk) String() string {
	return fmt.Sprintf("%d:%d:%s:%s:%d", c.Category, c.Depth, c.ParentID, c.ID, c.Tokens)
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
	dm := make(map[uint]bool)

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

func getChildChunk(chunks ChunkSlice, chunkID string) *Chunk {
	for _, chunk := range chunks {
		if chunk.ParentID == chunkID {
			if chunk.Category != CatText {
				return getChildChunk(chunks, chunk.ID)
			}
			return chunk
		}
	}
	return nil
}

func FindTextChunks(chunks ChunkSlice, summaryDepth uint) ChunkSlice {
	textChunks := make(ChunkSlice, 0, 21)
	for _, chunk := range chunks {
		if chunk.Category == CatSummary && chunk.Depth == summaryDepth {
			textChunk := getChildChunk(chunks, chunk.ID)
			textChunks = append(textChunks, textChunk)
		}
	}

	return textChunks
}
