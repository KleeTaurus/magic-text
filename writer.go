package magictext

import (
	"encoding/json"
	"fmt"
	"os"
)

// DumpChunksToJSON writes chunk slice to the given file.
func DumpChunksToJSON(filename string, chunks interface{}) error {
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

func DumpSummary(filename, summary string, captionSummaries []*CaptionSummary) error {
	file, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.WriteString(summary + "\n\n")
	if err != nil {
		return err
	}

	for _, cs := range captionSummaries {
		header := fmt.Sprintf("ID: %s, Seq: %d, Start: %s, End: %s", cs.ID, cs.Seq, cs.FromInString(), cs.ToInString())
		if _, err := file.WriteString(header + "\n"); err != nil {
			return err
		}

		if _, err := file.WriteString(cs.Text + "\n\n"); err != nil {
			return err
		}
	}

	return nil
}
