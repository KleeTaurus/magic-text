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

func SplitSubtitle(subtitle subtitles.Subtitle) ([]*CaptionChunk, error) {
	sb := strings.Builder{}
	chunks := make([]*CaptionChunk, 0, 10)

	var start time.Time
	seq := 0
	for i, caption := range subtitle.Captions {
		if i == 0 {
			start = caption.Start
		}

		text := strings.Join(caption.Text, " ")
		if CountTokens(sb.String()+text) > MaxReqTokens2048 {
			if i-1 < 0 {
				return nil, fmt.Errorf("Get caption end time failed")
			}
			chunks = append(chunks, NewCaptionChunk(seq, sb.String(), start, subtitle.Captions[i-1].End))

			sb.Reset()
			start = caption.Start
			seq++
		}
		if _, err := sb.WriteString(text + " "); err != nil {
			return nil, err
		}

		if i == len(subtitle.Captions)-1 && len(sb.String()) > 0 {
			chunks = append(chunks, NewCaptionChunk(seq, sb.String(), start, caption.End))
		}
	}

	return chunks, nil
}

func SplitText(text string, chunkSize, chunkOverlap int) ([]string, error) {
	if chunkSize > MaxReqTokens2048 {
		return nil, fmt.Errorf("The max tokens per chunk is %d, got %d", MaxReqTokens2048, chunkSize)
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
		if _, err := sb.WriteString(text + " "); err != nil {
			return nil, err
		}
	}

	if len(sb.String()) > 0 {
		chunks = append(chunks, strings.TrimSpace(sb.String())) // remove trailing spaces
	}

	return chunks, nil
}
