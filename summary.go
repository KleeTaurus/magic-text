package magictext

import (
	"fmt"
	"log"
	"sort"
	"sync"
	"time"
)

type Summary struct {
	ID   string
	Seq  int
	Text string
}

type SubtitleSummary struct {
	From time.Time
	Summary
}

const (
	MaxGroupChunks = 3
	MaxConcurrent  = 5
)

func recursiveSummary(topic string, chunks ChunkSlice, depth int) (ChunkSlice, error) {
	parentChunksMap := make(map[int]*Chunk)
	limiter := make(chan struct{}, MaxConcurrent)
	var wg sync.WaitGroup

	for i, chunkGroup := range chunks.SubGroups(MaxTokens2048, MaxGroupChunks) {
		limiter <- struct{}{}
		wg.Add(1)

		go func(chunkGroup ChunkSlice, i int, parentChunksMap map[int]*Chunk) {
			defer func() {
				<-limiter
				wg.Done()
			}()

			parentChunk := getParentChunk(i, topic, depth, chunkGroup)
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
		grandParentChunks, err := recursiveSummary(topic, parentChunks, depth+1)
		if err != nil {
			return nil, err
		}
		parentChunks = append(parentChunks, grandParentChunks...)
	}

	return parentChunks, nil
}

func getParentChunk(seq int, topic string, depth int, groupChunks ChunkSlice) *Chunk {
	log.Printf("%s, Generating text summary by openai.\n", groupChunks)
	var prompt string
	if topic != "" {
		prompt = fmt.Sprintf(GenerateSummaryPromptWithTopic, topic, groupChunks.TokenString())
	} else {
		prompt = fmt.Sprintf(GenerateSummaryPrompt, groupChunks.TokenString())
	}

	summary, err := completionWithRetry(prompt)
	if err != nil {
		log.Printf("%s, Generating text summary failed, err: %v", groupChunks, err)
		return NewSummaryChunk(seq, "", depth)
	}

	parentChunk := NewSummaryChunk(seq, summary, depth)
	// Update the child chunk's parent id
	for _, chunk := range groupChunks {
		chunk.ParentID = parentChunk.ID
	}

	return parentChunk
}
