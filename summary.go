package magictext

import (
	"fmt"
	"log"
	"sync"
)

const (
	MaxChunksPerGroup = 3 // TODO: this variable should be calculated dynamically
	MaxConcurrent     = 5
)

func generateSummary(topic string, chunks ChunkSlice) (*Chunk, error) {
	return summarizeRecursively(topic, chunks, 0)
}

func summarizeRecursively(topic string, chunks ChunkSlice, height int) (*Chunk, error) {
	summarizedChunksMap := make(map[int]*Chunk)

	mutex := &sync.Mutex{}
	limiter := make(chan struct{}, MaxConcurrent)
	var wg sync.WaitGroup

	chunkGroups := groupChunks(chunks, height)
	for i, chunkGroup := range chunkGroups {
		limiter <- struct{}{}
		wg.Add(1)

		go func(seq int, chunkGroup ChunkSlice) {
			defer func() {
				<-limiter
				wg.Done()
			}()

			summary, _ := summarizeChunks(topic, chunkGroup)
			summarizedChunk := NewChunk(seq, summary)
			summarizedChunk.Height = height + 1
			summarizedChunk.Children = chunkGroup

			mutex.Lock()
			summarizedChunksMap[seq] = summarizedChunk
			mutex.Unlock()
		}(i, chunkGroup)
	}
	wg.Wait()

	summarizedChunks := make(ChunkSlice, 0, len(summarizedChunksMap))
	for i := 0; i < len(summarizedChunksMap); i++ {
		summarizedChunks = append(summarizedChunks, summarizedChunksMap[i])
	}

	if len(summarizedChunks) == 1 {
		return summarizedChunks[0], nil
	}

	return summarizeRecursively(topic, summarizedChunks, height+1)
}

func summarizeChunks(topic string, groupChunks ChunkSlice) (string, error) {
	log.Printf("Generating text summary by openai, %s\n", groupChunks)
	var prompt string
	if topic != "" {
		prompt = fmt.Sprintf(GenerateSummaryPromptWithTopic, topic, groupChunks.Text())
	} else {
		prompt = fmt.Sprintf(GenerateSummaryPrompt, groupChunks.Text())
	}

	summary, err := completionWithRetry(prompt)
	if err != nil {
		log.Printf("Generating text summary failed, %s, err: %v", groupChunks, err)
		return "", err
	}

	return summary, nil
}

func groupChunks(cs ChunkSlice, height int) []ChunkSlice {
	// height 0, max chunks = 1
	// height 1, max chunks = 4
	// height 2, max chunks = 7
	// height 3, max chunks = 10
	maxChunksInGroup := MaxChunksPerGroup*height + 1
	groups := make([]ChunkSlice, 0, len(cs)/2)
	chunkGroup := make(ChunkSlice, 0, maxChunksInGroup)
	for _, chunk := range cs {
		if chunkGroup.Tokens()+chunk.Tokens > MaxReqTokens2048 || len(chunkGroup) >= maxChunksInGroup {
			groups = append(groups, chunkGroup)
			// reset chunkGroup to empty
			chunkGroup = make(ChunkSlice, 0, maxChunksInGroup)
		}
		chunkGroup = append(chunkGroup, chunk)
	}
	groups = append(groups, chunkGroup)
	return groups
}
