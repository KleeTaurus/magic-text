package summaryit

import (
	"log"
	"strings"
	"unicode/utf8"
)

const (
	NoLimitOnDepth = 18446744073709551615
	MaxGroupChunks = 4
)

func SummaryFile(prompt, filename string) (string, error) {
	textChunks := ReadTextFile(filename)

	log.Println("Total text chunks: ", len(textChunks))
	summaryChunks, err := recursiveSummary(prompt, 0, textChunks)
	if err != nil {
		return "", err
	}

	summaryFile := getOutfile(filename, ".sum")
	WriteTextFile(summaryFile, summaryChunks)
	WriteJSONFile(getOutfile(filename, ".json"), append(textChunks, summaryChunks...))

	return summaryFile, nil
}

func recursiveSummary(prompt string, depth uint, chunks ChunkSlice) (ChunkSlice, error) {
	parentChunks := make(ChunkSlice, 0, len(chunks)/2)

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
		grandParentChunks, err := recursiveSummary(prompt, depth+1, parentChunks)
		if err != nil {
			return parentChunks, err
		}
		parentChunks = append(parentChunks, grandParentChunks...)
	}

	return parentChunks, nil
}
