package magictext

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
)

// DumpChunkToJSON writes chunk slice to the given file.
func DumpChunkToJSON(filename string, chunks ChunkSlice) error {
	file, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	b, err := json.MarshalIndent(chunks, "", "  ")
	if err != nil {
		return err
	}

	if _, err := file.Write(b); err != nil {
		return err
	}

	return nil
}

func DumpChunkToText(filename string, chunks ChunkSlice) error {
	file, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	var preLevel int = 0
	for i, chunk := range chunks {
		if i == 0 || chunk.Depth != preLevel {
			text := "# LEVEL" + strconv.Itoa(int(chunk.Depth))
			if _, err := file.WriteString(text + "\n\n"); err != nil {
				return err
			}
			preLevel = chunk.Depth
		}

		if _, err := file.WriteString(fmt.Sprintf("%d %s", chunk.Seq, chunk.Text) + "\n\n"); err != nil {
			return err
		}
	}

	return nil
}
