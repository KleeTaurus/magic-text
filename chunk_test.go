package summaryit

import (
	"fmt"
	"testing"
)

func TestGroup(t *testing.T) {
	chunks := make(ChunkSlice, 0, 10)
	for i := 0; i < 10; i++ {
		chunks = append(chunks, NewTextChunk(fmt.Sprintf("This is a chunk text %03d", i))) // run count: 24
	}

	groups := chunks.SubGroups(48, 3)
	if len(groups) != 5 {
		t.Errorf("expected 5 groups, got %d", len(groups))
	}
	for _, g := range groups {
		if g.Tokens() != 48 {
			t.Errorf("expected 48 tokens, got %d", g.Tokens())
		}
	}

	groups = chunks.SubGroups(250, 4)
	if len(groups) != 3 {
		t.Errorf("expected 3 groups, got %d", len(groups))
	}
}

func TestCountTokens(t *testing.T) {
	text1 := "你好，世界！"
	text2 := "Hello, World!"
	tokenCount := CountTokens(text1)
	if tokenCount != 7 {
		t.Errorf("expected 7 tokens, got %d", tokenCount)
	}

	tokenCount = CountTokens(text2)
	if tokenCount != 4 {
		t.Errorf("expected 4 tokens, got %d", tokenCount)
	}
}
