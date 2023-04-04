package textsummary

import (
	"context"
	"log"
	"os"

	"github.com/sashabaranov/go-openai"
)

func summaryByOpenAI(prompt, content string) string {
	log.Println("Generating text summary by openai...")
	log.Println("Prompt: ", prompt, " content: ", content)

	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		log.Fatal("Missing OpenAI API key")
	}

	client := openai.NewClient(apiKey)
	resp, err := client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: openai.GPT3Dot5Turbo,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleUser,
					Content: prompt + " " + content,
				},
			},
		},
	)

	if err != nil {
		log.Printf("ChatCompletion error: %v\n", err)
		return ""
	}

	result := resp.Choices[0].Message.Content
	log.Println("ChatCompletion result: ", result)
	log.Println("The text summary has been successfully generated")
	return result
}
