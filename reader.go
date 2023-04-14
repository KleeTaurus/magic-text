package summaryit

import (
	"bufio"
	"encoding/json"
	"log"
	"os"
	"strings"
)

// ReadTextFile reads a text file and stores it's content in the Chunk slice.
func ReadTextFile(filename string) ChunkSlice {
	file, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	chunks := make(ChunkSlice, 0, 10)
	sb := strings.Builder{}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		// ignore empty lines
		if strings.TrimSpace(line) == "" {
			continue
		}

		if CountTokens(sb.String())+CountTokens(line) > MaxTokensPerRequest {
			chunks = append(chunks, NewTextChunk(sb.String()))
			sb.Reset()
		}

		// Here we separate the paragraph by a space
		sb.WriteString(line + " ")
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	// Append the remaining text to the chunk slice
	chunks = append(chunks, NewTextChunk(sb.String()))

	return chunks
}

func ReadJSONFile(filename string) ChunkSlice {
	data, err := os.ReadFile(filename)
	if err != nil {
		log.Fatal(err)
	}

	var chunks ChunkSlice
	if err := json.Unmarshal(data, &chunks); err != nil {
		log.Fatal(err)
	}
	return chunks
}
