package magictext

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/martinlindhe/subtitles"
)

var (
	re = regexp.MustCompile("\\n+")
)

func SplitSubtitle(subtitle subtitles.Subtitle) ([]CaptionChunk, error) {
	sb := strings.Builder{}
	chunks := make([]CaptionChunk, 0, 21)

	var start time.Time
	for i, caption := range subtitle.Captions {
		if i == 0 {
			start = caption.Start
		}

		text := strings.Join(caption.Text, " ")
		if CountTokens(sb.String()+text) > MaxTokens2048 {
			chunks = append(chunks, NewCaptionChunk(start, sb.String()))

			sb.Reset()
			start = caption.Start
		}
		sb.WriteString(text + " ")

		if i == len(subtitle.Captions)-1 && len(sb.String()) > 0 {
			chunks = append(chunks, NewCaptionChunk(start, sb.String()))
		}
	}

	return chunks, nil
}

func SplitText(text string, chunkSize, chunkOverlap int) ([]string, error) {
	if chunkSize > MaxTokens2048 {
		return nil, fmt.Errorf("The max tokens per chunk is %d, got %d", MaxTokens2048, chunkSize)
	}

	if chunkOverlap > chunkSize {
		return nil, fmt.Errorf("Chunk overlap %d is larger than chunk size %d", chunkOverlap, chunkSize)
	}

	tokens := CountTokens(text)
	if tokens < chunkSize {
		return []string{text}, nil
	}

	texts := re.Split(text, -1)
	chunks := make([]string, 0, tokens/chunkSize)
	sb := strings.Builder{}

	for _, text := range texts {
		if strings.TrimSpace(text) == "" {
			continue
		}

		if CountTokens(text) > chunkSize {
			return nil, fmt.Errorf("The length of text exceeds %d", chunkSize)
		}

		if CountTokens(sb.String()+text) > chunkSize {
			chunks = append(chunks, strings.TrimSpace(sb.String())) // remove trailing spaces
			sb.Reset()
		}
		sb.WriteString(text + " ")
	}

	if len(sb.String()) > 0 {
		chunks = append(chunks, strings.TrimSpace(sb.String())) // remove trailing spaces
	}

	return chunks, nil
}
