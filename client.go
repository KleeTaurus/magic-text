package magictext

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/martinlindhe/subtitles"
	"github.com/pkoukk/tiktoken-go"
	"github.com/sashabaranov/go-openai"
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
	MaxTokens512  = 512
	MaxTokens2048 = 2048
)

var (
	Debug        = false
	MockOpenAI   = false
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
func GenerateSummaryBySubtitle(topic string, subtitle subtitles.Subtitle) ([]*SubtitleSummary, string, error) {
	subtitleSummaries := make([]*SubtitleSummary, 0, 11)

	// Split subtitle into caption chunks
	captionChunks, err := SplitSubtitle(subtitle)
	if err != nil {
		return subtitleSummaries, "", err
	}

	// Save caption chunks into a map, so we can get start time
	// by content hash id
	captionChunksMap := make(map[string]CaptionChunk, 0)
	chunks := make(ChunkSlice, 0, len(captionChunks))
	for i, cc := range captionChunks {
		chunks = append(chunks, NewChunk(i, cc.Text))
		captionChunksMap[cc.ID] = cc
	}

	log.Println("Total chunks: ", len(chunks))
	rootChunk, err := generateSummary(topic, chunks)
	if err != nil {
		return subtitleSummaries, "", err
	}

	randomFile := randomFilename()
	DumpChunksToJSON("/tmp/"+randomFile+"_1.json", captionChunks)
	DumpChunksToJSON("/tmp/"+randomFile+"_2.json", chunks)
	DumpChunksToJSON("/tmp/"+randomFile+"_3.json", rootChunk)

	summary := rootChunk.Text
	for _, child := range rootChunk.Children {
		for _, grandchild := range child.Children {
			ss := &SubtitleSummary{}
			ss.ID = grandchild.ID
			ss.Seq = grandchild.Seq
			ss.Text = grandchild.Text

			leaf := getLeafChunk(grandchild)
			if cc, ok := captionChunksMap[leaf.ID]; ok {
				ss.From = cc.From
			}

			subtitleSummaries = append(subtitleSummaries, ss)
		}
	}

	return subtitleSummaries, summary, nil
}

func getLeafChunk(target *Chunk) *Chunk {
	if len(target.Children) == 0 {
		return target
	}
	return getLeafChunk(target.Children[0])
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
