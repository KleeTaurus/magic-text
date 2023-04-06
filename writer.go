package summaryit

import (
	"encoding/json"
	"log"
	"os"
	"strconv"
)

// WriteJSONFile writes Chunk slice to the given file.
func WriteJSONFile(filename string, chunks ChunkSlice) {
	file, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	b, err := json.MarshalIndent(chunks, "", "  ")
	if err != nil {
		log.Fatal(err)
	}

	if _, err := file.Write(b); err != nil {
		log.Fatal(err)
	}
}

func WriteTextFile(filename string, chunks ChunkSlice) {
	file, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	var preLevel uint = 0
	for i, chunk := range chunks {
		if i == 0 || chunk.Depth != preLevel {
			file.WriteString("# LEVEL" + strconv.Itoa(int(chunk.Depth)) + "\n\n")
			preLevel = chunk.Depth
		}

		if _, err := file.WriteString(chunk.Text + "\n\n"); err != nil {
			log.Fatal(err)
		}
	}
}
