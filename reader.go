package magictext

import (
	"bufio"
	"crypto/tls"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/martinlindhe/subtitles"
)

type Lines struct {
	text []string
}

func ReadSRTURI(uri string) ChunkSlice {
	// Create a http.Transport object
	transport := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	// Create http.Client object
	client := &http.Client{Transport: transport}

	// Send HTTPS request and get the response
	resp, err := client.Get(uri)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	// Read response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	return srtBytesToChunks(body)
}

func ReadSRTFile(filename string) ChunkSlice {
	b, err := os.ReadFile(filename)
	if err != nil {
		log.Fatal(err)
	}

	return srtBytesToChunks(b)
}

// ReadTextFile reads a text file and stores it's content in the Chunk slice.
func ReadTextFile(filename string) ChunkSlice {
	file, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	var lines Lines
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines.text = append(lines.text, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	return linesToChunks(lines)
}

func ReadJSONFile(filename string) ChunkSlice {
	b, err := os.ReadFile(filename)
	if err != nil {
		log.Fatal(err)
	}

	var chunks ChunkSlice
	if err := json.Unmarshal(b, &chunks); err != nil {
		log.Fatal(err)
	}
	return chunks
}

func srtBytesToChunks(b []byte) ChunkSlice {
	subtitles, err := subtitles.NewFromSRT(string(b))
	if err != nil {
		log.Fatal(err)
	}

	var lines Lines
	for _, caption := range subtitles.Captions {
		lines.text = append(lines.text, caption.Text...)
	}

	return linesToChunks(lines)
}

func linesToChunks(lines Lines) ChunkSlice {
	chunks := make(ChunkSlice, 0, 30)
	sb := strings.Builder{}

	for i, line := range lines.text {
		// ignore empty lines
		if strings.TrimSpace(line) == "" {
			continue
		}

		if CountTokens(sb.String())+CountTokens(line) > MaxTokens2048 {
			chunks = append(chunks, NewTextChunk(i, sb.String()))
			sb.Reset()
		}

		// Here we separate the paragraph by a space
		sb.WriteString(line + " ")

		if i == len(lines.text)-1 {
			chunks = append(chunks, NewTextChunk(i, sb.String()))
		}
	}

	return chunks
}
