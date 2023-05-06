package magictext

import (
	"fmt"
	"strings"
)

type Chunk struct {
	ID       string     `json:"id"`
	Height   int        `json:"height"`
	Seq      int        `json:"seq"`
	Text     string     `json:"text"`
	Tokens   int        `json:"tokens"`
	Children ChunkSlice `json:"children"`
}

func NewChunk(seq int, text string) *Chunk {
	return &Chunk{
		ID:       hashString(text),
		Height:   0,
		Seq:      seq,
		Text:     text,
		Tokens:   CountTokens(text),
		Children: ChunkSlice{},
	}
}

func (c *Chunk) String() string {
	return fmt.Sprintf("%s:%d:%d:%d:%d", c.ID, c.Height, c.Seq, c.Tokens, len(c.Children))
}

type ChunkSlice []*Chunk

func (cs ChunkSlice) Tokens() int {
	tokens := 0
	for _, chunk := range cs {
		tokens += chunk.Tokens
	}
	return tokens
}

func (cs ChunkSlice) Text() string {
	sb := strings.Builder{}
	for _, chunk := range cs {
		sb.WriteString(chunk.Text + " ") // TODO: replace the space with \n\n?
	}
	return sb.String()
}

func (cs ChunkSlice) String() string {
	heightMap := make(map[int]bool)
	heights := []string{}
	seqs := []string{}

	for _, chunk := range cs {
		seqs = append(seqs, fmt.Sprintf("%02d", chunk.Seq))
		if _, ok := heightMap[chunk.Height]; !ok {
			heightMap[chunk.Height] = true
			heights = append(heights, fmt.Sprintf("%d", chunk.Height))
		}
	}

	return fmt.Sprintf("Heights: %s, Seqs: %s, Chunks: %d, Tokens: %d",
		strings.Join(heights, "_"), strings.Join(seqs, "_"), len(cs), cs.Tokens())
}
