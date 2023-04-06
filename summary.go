package summaryit

import (
	"fmt"
	"log"
	"path"
	"strings"
	"unicode/utf8"
)

const (
	NoLimitOnDepth = 18446744073709551615
)

func SummaryFile(prompt, filename string, maxDepth uint) (string, error) {
	textChunks := ReadTextFile(filename)

	log.Println("Total text chunks: ", len(textChunks))
	summaryChunks, err := recursiveSummary(prompt, 0, maxDepth, textChunks)
	if err != nil {
		return "", err
	}

	basename := path.Base(filename)
	outfile := fmt.Sprintf("%s/%s.json", path.Dir(filename), strings.Split(basename, ".")[0])
	WriteJSONFile(outfile, append(textChunks, summaryChunks...))

	return outfile, err
}

func recursiveSummary(prompt string, depth uint, maxDepth uint, chunks ChunkSlice) (ChunkSlice, error) {
	parentChunks := make(ChunkSlice, 0, len(chunks)/2)
	if depth > maxDepth {
		log.Printf("Exceeding max depth, current depth: %d, max depth: %d", depth, maxDepth)
		return parentChunks, nil
	}

	tokens := strings.Builder{}
	start := 0
	for i, chunk := range chunks {
		if utf8.RuneCountInString(tokens.String())+chunk.RuneCountInString() > maxTokens {
			summary := summaryByOpenAI(prompt, tokens.String())
			parentChunk := NewSummaryChunk(summary, depth)
			parentChunks = append(parentChunks, parentChunk)

			for j := start; j < i; j++ {
				chunks[j].ParentID = parentChunk.ID
			}

			start = i
			tokens.Reset()
		}
		tokens.WriteString(chunk.Text)
	}

	summary := summaryByOpenAI(prompt, tokens.String())
	parentChunk := NewSummaryChunk(summary, depth)
	parentChunks = append(parentChunks, parentChunk)

	for j := start; j < len(chunks); j++ {
		chunks[j].ParentID = parentChunk.ID
	}

	if len(parentChunks) > 1 {
		grandParentChunks, err := recursiveSummary(prompt, depth+1, maxDepth, parentChunks)
		if err != nil {
			return parentChunks, err
		}
		parentChunks = append(parentChunks, grandParentChunks...)
	}

	return parentChunks, nil
}
