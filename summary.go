package summaryit

import (
	"fmt"
	"log"
	"sort"
	"sync"
)

const (
	MaxGroupChunks      = 2
	MaxConcurrent       = 3
	BaseChunkDepth      = 0
	MaxTokensPerRequest = 1024
)

func SummaryFile(prompt, filename string) (string, error) {
	textChunks := ReadTextFile(filename)

	fmt.Println("Total chunks of input file: ", len(textChunks))
	summaryChunks, err := RecursiveSummary(prompt, textChunks, BaseChunkDepth)
	if err != nil {
		return "", err
	}

	summaryFile := GenerateFilename(filename, ".sum")
	WriteTextFile(summaryFile, summaryChunks)
	WriteJSONFile(GenerateFilename(filename, ".json"), append(textChunks, summaryChunks...))

	return summaryFile, nil
}

func RecursiveSummary(prompt string, chunks ChunkSlice, depth uint) (ChunkSlice, error) {
	parentChunksMap := make(map[int]*Chunk)
	limiter := make(chan struct{}, MaxConcurrent)
	var wg sync.WaitGroup

	for i, chunkGroup := range chunks.SubGroups(MaxTokensPerRequest, MaxGroupChunks) {
		limiter <- struct{}{}
		wg.Add(1)

		go func(chunkGroup ChunkSlice, i int, parentChunksMap map[int]*Chunk) {
			defer func() {
				<-limiter
				wg.Done()
			}()

			parentChunk := getParentChunk(prompt, depth, chunkGroup)
			parentChunksMap[i] = parentChunk
		}(chunkGroup, i, parentChunksMap)
	}
	wg.Wait()

	keys := make([]int, 0, len(parentChunksMap))
	for key := range parentChunksMap {
		keys = append(keys, key)
	}
	sort.Ints(keys)

	parentChunks := make(ChunkSlice, 0, len(keys))
	for key := range keys {
		parentChunks = append(parentChunks, parentChunksMap[key])
	}

	if len(parentChunks) > 1 {
		grandParentChunks, err := RecursiveSummary(prompt, parentChunks, depth+1)
		if err != nil {
			return nil, err
		}
		parentChunks = append(parentChunks, grandParentChunks...)
	}

	return parentChunks, nil
}

func getParentChunk(prompt string, depth uint, groupChunks ChunkSlice) *Chunk {
	log.Printf("%s, Generating text summary by openai.\n", groupChunks)
	summary := summaryByOpenAI(prompt, groupChunks.TokenString())
	log.Printf("%s, The text summary has been successfully generated.\n", groupChunks)

	parentChunk := NewSummaryChunk(summary, depth)
	// Update the child chunk's parent id
	for _, chunk := range groupChunks {
		chunk.ParentID = parentChunk.ID
	}

	return parentChunk
}
