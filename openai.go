package magictext

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/sashabaranov/go-openai"
)

const (
	MaxRetryTimes = 3
	SleepSeconds  = 3
)

func completionWithRetry(prompt string) (string, error) {
	return retry(completion, prompt, MaxRetryTimes)
}

func retry(fn func(string) (string, error), prompt string, retryTimes int) (string, error) {
	for i := 1; i <= retryTimes; i++ {
		str, err := fn(prompt)
		if err == nil {
			return str, nil
		}

		log.Printf("%d: ChatCompletion error: %v, retry after %d seconds\n", i, err, SleepSeconds)
		if i != retryTimes {
			time.Sleep(time.Second * SleepSeconds)
		}
	}
	return "", fmt.Errorf("retry failed after %d times", retryTimes)
}

func completion(prompt string) (string, error) {
	if Debug {
		log.Printf("prompt:\n%s", prompt)
	}

	if MockOpenAI {
		return fmt.Sprintf("It's a mock response of openai api, data: %s", randString()), nil
	}

	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		return "", fmt.Errorf("Missing OpenAI API key, You must set OPENAI_API_KEY environment variable first")
	}

	OpenAIClient = openai.NewClient(apiKey)
	resp, err := OpenAIClient.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model:       openai.GPT3Dot5Turbo,
			MaxTokens:   500,
			Temperature: 0.2,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleUser,
					Content: prompt,
				},
			},
		},
	)

	if err != nil {
		return "", err
	}

	if Debug {
		log.Printf("response: %+v\n", resp.Choices)
	}

	return resp.Choices[0].Message.Content, nil
}
