package magictext

import (
	"fmt"
	"reflect"
	"testing"
)

func TestChunkSlice_SubGroups(t *testing.T) {
	size := 9
	chunks := make(ChunkSlice, 0, size)
	for i := 0; i < size; i++ {
		chunks = append(chunks, NewTextChunk(fmt.Sprintf("This is a chunk %d", i)))
	}

	type args struct {
		maxTokensPerRequest int
		maxChunksInGroup    int
	}
	tests := []struct {
		name string
		cs   ChunkSlice
		args args
		want []ChunkSlice
	}{
		{"Divide by max chunks in group", chunks, args{30, 3}, []ChunkSlice{
			{chunks[0], chunks[1], chunks[2]},
			{chunks[3], chunks[4], chunks[5]},
			{chunks[6], chunks[7], chunks[8]},
		}},
		{"Divide by max tokens", chunks, args{21, 5}, []ChunkSlice{
			{chunks[0], chunks[1], chunks[2]},
			{chunks[3], chunks[4], chunks[5]},
			{chunks[6], chunks[7], chunks[8]},
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.cs.SubGroups(tt.args.maxTokensPerRequest, tt.args.maxChunksInGroup); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ChunkSlice.SubGroups() = %v, want %v", got, tt.want)
			}
		})
	}
}
