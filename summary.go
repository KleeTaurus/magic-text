package summaryit

import (
	"fmt"
	"log"
	"sort"
	"sync"
	"time"
)

const (
	MaxGroupChunks      = 3
	MaxConcurrent       = 3
	BaseChunkDepth      = 0
	MaxTokensPerRequest = 2048
)

var (
	prompt = "`reset` " +
		"`no quotes` " +
		"`no explanations` " +
		"`no prompt` " +
		"`no self-reference` " +
		"`no apologies` " +
		"`no filler` " +
		"`just answer` " + `
I will give you text content, you will rewrite it and output that in a short 
summarized version of my text. Keep the meaning the same. Ensure that the 
revised content has significantly fewer characters than the original text, 
and no more than 500 Chinese words, the fewer the better.
` +
		`When generating text summaries, expand around the following topics as
much as possible: %s` +
		`Only give me the output and nothing else. Now, using the concepts above, 
summarize the following text. Respond in Chinese language:
`
)

func SummaryFile(customPrompt, filename string) (string, error) {
	textChunks := ReadSRTFile(filename)

	fmt.Println("Total chunks of input file: ", len(textChunks))
	summaryChunks, err := RecursiveSummary(fmt.Sprintf(prompt, customPrompt), textChunks, BaseChunkDepth)
	if err != nil {
		return "", err
	}

	summaryFile := MakeFilename(filename, "sum")
	jsonFile := MakeFilename(filename, "json")
	WriteTextFile(summaryFile, summaryChunks)
	WriteJSONFile(jsonFile, append(textChunks, summaryChunks...))

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
	summary, err := retry(summaryByOpenAI, prompt, groupChunks.TokenString(), 3)
	if err != nil {
		log.Printf("%s, Generating text summary failed, err: %v", groupChunks, err)
		return NewSummaryChunk("", depth)
	}
	// log.Printf("%s, The text summary has been successfully generated.\n", groupChunks)

	parentChunk := NewSummaryChunk(summary, depth)
	// Update the child chunk's parent id
	for _, chunk := range groupChunks {
		chunk.ParentID = parentChunk.ID
	}

	return parentChunk
}

func retry(fn func(string, string) (string, error), prompt, content string, times int) (string, error) {
	for i := 0; i < times; i++ {
		str, err := fn(prompt, content)
		if err == nil {
			return str, nil
		}
		log.Printf("[%d] Calling OpenAI API failed, err: %v", i, err)
		time.Sleep(time.Second * 5)
	}
	return "", fmt.Errorf("retry failed for %d times", times)
}
