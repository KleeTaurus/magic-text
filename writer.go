package summaryit

import (
	"encoding/json"
	"log"
	"os"
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
