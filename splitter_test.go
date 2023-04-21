package magictext

import (
	"testing"
)

const Paragraph = `
This is an era where AI breakthrough is coming daily. We didn’t have many AI-generated in public a few years ago, but now the technology is accessible to everyone. It’s excellent for many individual creators or companies that want to significantly take advantage of the technology to develop something complex, which might take a long time.

One of the most incredible breakthroughs that change how we work is the release of the GPT-3.5 model by OpenAI. What is the GPT-3.5 model? If I let the model talk for themselves. In that case, the answer is “a highly advanced AI model in the field of natural language processing, with vast improvements in generating contextually accurate and relevant text”.

OpenAI provides an API for the GPT-3.5 model that we can use to develop a simple app, such as a text summarizer. To do that, we can use Python to integrate the model API into our intended application seamlessly. What does the process look like? Let’s get into it.
`

func TestSplit(t *testing.T) {
	type args struct {
		text         string
		chunkSize    int
		chunkOverlap int
	}
	tests := []struct {
		name    string
		args    args
		want    int
		wantErr bool
	}{
		{"Split long text 140 tokens per chunk", args{Paragraph, 140, 0}, 3, false},
		{"Split long text 30 tokens per chunk", args{Paragraph, 30, 0}, 0, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := SplitText(tt.args.text, tt.args.chunkSize, tt.args.chunkOverlap)
			if (err != nil) != tt.wantErr {
				t.Errorf("Split() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if len(got) != tt.want {
				t.Errorf("Split() = %v, want %v", len(got), tt.want)
			}

			/*
				for _, chunk := range got {
					fmt.Printf("[%d] %s\n", CountTokens(chunk), chunk)
				}
			*/
		})
	}
}
