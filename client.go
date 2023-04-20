package magictext

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/pkoukk/tiktoken-go"
	"github.com/sashabaranov/go-openai"
)

const (
	MaxTokensInGenerateTitle = 512
	MaxTokensInExtractNouns  = 2048
)

var (
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

// GenerateSummary generates a summary for the given text
func GenerateSummary(longtext string, topic string) []Summary {
	return []Summary{}
}

// GenerateTitle generates a title for the given text, the max length of input text is 512.
func GenerateTitle(text string) (string, error) {
	if tokens, ok := validateTokens(text, MaxTokensInGenerateTitle); !ok {
		return "", fmt.Errorf("The maximum tokens supported is %d, got %d", MaxTokensInGenerateTitle, tokens)
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
	if tokens, ok := validateTokens(text, MaxTokensInExtractNouns); !ok {
		return "", fmt.Errorf("The maximum tokens supported is %d, got %d", MaxTokensInExtractNouns, tokens)
	}

	jsonStr, err := completionWithRetry(fmt.Sprintf(ExtractNounsPrompt, text))
	if err != nil {
		return "", err
	}

	return jsonStr, nil
}
