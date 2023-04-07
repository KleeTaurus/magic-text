package summaryit

import (
	"fmt"
	"log"
)

const (
	MaxGroupChunks     = 2
	MaxCurrentRequests = 3
	BaseChunkDepth     = 0
)

func SummaryFile(prompt, filename string) (string, error) {
	textChunks := ReadTextFile(filename)

	fmt.Println("Total chunks of input file: ", len(textChunks))
	summaryChunks, err := RecursiveSummary(prompt, textChunks, BaseChunkDepth)
	if err != nil {
		return "", err
	}

	summaryFile := getOutfile(filename, ".sum")
	WriteTextFile(summaryFile, summaryChunks)
	WriteJSONFile(getOutfile(filename, ".json"), append(textChunks, summaryChunks...))

	return summaryFile, nil
}

func RecursiveSummary(prompt string, chunks ChunkSlice, depth uint) (ChunkSlice, error) {
	parentChunks := make(ChunkSlice, 0, len(chunks)/2)
	childChunks := make(ChunkSlice, 0, len(chunks))

	for _, chunk := range chunks {
		if childChunks.Tokens()+chunk.Tokens() > maxTokens ||
			len(childChunks) >= MaxGroupChunks {
			parentChunks = addParentChunk(prompt, depth, parentChunks, childChunks)

			// Clear the child chunks
			childChunks = make(ChunkSlice, 0, len(chunks))
		}
		childChunks = append(childChunks, chunk)
	}
	parentChunks = addParentChunk(prompt, depth, parentChunks, childChunks)

	if len(parentChunks) > 1 {
		grandParentChunks, err := RecursiveSummary(prompt, parentChunks, depth+1)
		if err != nil {
			return nil, err
		}
		parentChunks = append(parentChunks, grandParentChunks...)
	}

	return parentChunks, nil
}

func addParentChunk(prompt string, depth uint, parentChunks, childChunks ChunkSlice) ChunkSlice {
	log.Printf("%s, Generating text summary by openai.\n", childChunks)
	summary := summaryByOpenAI(prompt, childChunks.TokenString())
	log.Printf("%s, The text summary has been successfully generated.\n", childChunks)

	parentChunk := NewSummaryChunk(summary, depth)
	// Update the child chunk's parent id
	for _, chunk := range childChunks {
		chunk.ParentID = parentChunk.ID
	}

	parentChunks = append(parentChunks, parentChunk)
	return parentChunks
}
