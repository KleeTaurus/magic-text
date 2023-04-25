package magictext

import (
	"fmt"
	"log"
)

const (
	MaxChunksPerGroup = 3 // TODO: this variable should be calculated dynamically
	MaxConcurrent     = 5
	BaseHeight        = 0
)

func generateSummary(topic string, chunks ChunkSlice) (*Chunk, error) {
	return summarizeRecursively(topic, chunks, BaseHeight)
}

func summarizeRecursively(topic string, chunks ChunkSlice, height int) (*Chunk, error) {
	summarizedChunksMap := make(map[int]*Chunk)

	chunkGroups := groupChunks(chunks, MaxTokens2048, height*MaxChunksPerGroup)
	for i, chunkGroup := range chunkGroups {
		func(seq int, chunkGroup ChunkSlice) {
			summary, _ := summarizeChunks(topic, chunkGroup)
			summarizedChunk := NewChunk(seq, summary)
			summarizedChunk.Height = height + 1
			summarizedChunk.Children = chunkGroup

			summarizedChunksMap[i] = summarizedChunk
		}(i, chunkGroup)
	}

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
		return "", err
	}

	return summary, nil
}

func groupChunks(cs ChunkSlice, maxTokensPerRequest, maxChunksInGroup int) []ChunkSlice {
	groups := make([]ChunkSlice, 0, len(cs)/2)
	chunkGroup := make(ChunkSlice, 0, maxChunksInGroup)
	for _, chunk := range cs {
		if chunkGroup.Tokens()+chunk.Tokens > maxTokensPerRequest || len(chunkGroup) > maxChunksInGroup {
			groups = append(groups, chunkGroup)
			// reset chunkGroup to empty
			chunkGroup = make(ChunkSlice, 0, maxChunksInGroup)
		}
		chunkGroup = append(chunkGroup, chunk)
	}
	groups = append(groups, chunkGroup)
	return groups
}

func calculateMaxChunksInGroup(cs ChunkSlice) int {
	if len(cs) < 11 {
		return 11
	}

	return 11
}
