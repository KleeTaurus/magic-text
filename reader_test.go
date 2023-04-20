package magictext

import (
	"reflect"
	"testing"
)

func TestReadSRTFile(t *testing.T) {
	type args struct {
		filename string
	}
	tests := []struct {
		name string
		args args
		want ChunkSlice
	}{
		{"", args{"example/data/The.Godfather.I.srt"}, make(ChunkSlice, 0, 0)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ReadSRTFile(tt.args.filename); reflect.DeepEqual(got, tt.want) {
				t.Errorf("ReadSRTFile() = %v, want %v", got, tt.want)
			}
		})
	}
}
