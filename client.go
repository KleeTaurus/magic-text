package magictext

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/martinlindhe/subtitles"
	"github.com/pkoukk/tiktoken-go"
	"github.com/sashabaranov/go-openai"
)

const (
	MaxTokens512  = 512
	MaxTokens2048 = 2048
)

var (
	Debug        = false
	OpenAIClient *openai.Client
	TikToken     *tiktoken.Tiktoken
)

func init() {
	// 1. get openai api key
	// 2. initialize tiktoken
	godotenv.Load()

	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		log.Fatal("Missing OpenAI API key, You must set OPENAI_API_KEY environment variable first")
	}
	OpenAIClient = openai.NewClient(apiKey)

	var err error
	TikToken, err = tiktoken.GetEncoding("cl100k_base") // support models: gpt-4, gpt-3.5-turbo, text-embedding-ada-002
	if err != nil {
		log.Fatal(err)
	}
}

// GenerateSummaryBySubtitle generates a summary for the given subtitles
func GenerateSummaryBySubtitle(topic string, subtitle subtitles.Subtitle) ([]SubtitleSummary, string, error) {
	subtitleSummaries := make([]SubtitleSummary, 0, 11)

	// Split subtitle into caption chunks
	captionChunks, err := SplitSubtitle(subtitle)
	if err != nil {
		return subtitleSummaries, "", err
	}

	randomFile := randomFilename()
	captionChunkFile := "/tmp/" + randomFile + ".1.json"
	DumpChunksToJSON(captionChunkFile, captionChunks)

	captionMap := make(map[string]CaptionChunk, 0)
	textChunks := make(ChunkSlice, 0, len(captionChunks))
	for i, cc := range captionChunks {
		textChunks = append(textChunks, NewTextChunk(i, cc.Text))
		captionMap[cc.ID] = cc
	}

	textChunkFile := "/tmp/" + randomFile + ".2.json"
	DumpChunksToJSON(textChunkFile, textChunks)

	fmt.Println("Total text chunks: ", len(textChunks))
	summaryChunks, err := recursiveSummary(topic, textChunks, 0)
	if err != nil {
		return subtitleSummaries, "", err
	}

	summaryChunkFile := "/tmp/" + randomFile + ".3.json"
	DumpChunksToJSON(summaryChunkFile, summaryChunks)

	var summary string
	for i, sumChunk := range summaryChunks {
		if i == len(summaryChunks)-1 {
			summary = sumChunk.Text
		}

		// we only want level 1 chunks
		if sumChunk.Depth != 1 {
			continue
		}

		ps := SubtitleSummary{}
		ps.ID = sumChunk.ID
		ps.Seq = sumChunk.Seq
		ps.Text = sumChunk.Text

		for _, c := range summaryChunks {
			if ps.ID == c.ParentID {
				for _, t := range textChunks {
					if c.ID == t.ParentID {
						if cc, ok := captionMap[t.ID]; ok {
							ps.From = cc.From
						}
					}
				}
			}
		}
		subtitleSummaries = append(subtitleSummaries, ps)
	}

	return subtitleSummaries, summary, nil
}

// GenerateTitle generates a title for the given text, the max length of input text is 512.
func GenerateTitle(text string) (string, error) {
	if tokens, ok := validateTokens(text, MaxTokens512); !ok {
		return "", fmt.Errorf("The maximum tokens supported is %d, got %d", MaxTokens512, tokens)
	}

	result, err := completionWithRetry(fmt.Sprintf(GenerateTitlePrompt, text))
	if err != nil {
		return "", err
	}

	return result, nil
}

// ExtractNouns extracts nouns from a string, the max length of input text is 2048, the output
// is a json string, see following example for more information.
//
// Output string:
//
//	{
//	   "usernames": ["吴三桂", "皇太极", "弘历"],
//	   "company_names": ["得到"],
//	   "product_names": [],
//	   "course_names": ["硅谷来信"],
//	   "book_names": ["万历十五年", "湘行散记", "货币未来"]
//	}
func ExtractNouns(text string) (string, error) {
	if tokens, ok := validateTokens(text, MaxTokens2048); !ok {
		return "", fmt.Errorf("The maximum tokens supported is %d, got %d", MaxTokens2048, tokens)
	}

	result, err := completionWithRetry(fmt.Sprintf(ExtractNounsPrompt, text))
	if err != nil {
		return "", err
	}

	return result, nil
}
