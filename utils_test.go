package summaryit

import (
	"testing"
)

func TestCountTokens(t *testing.T) {
	type args struct {
		text string
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{"Count English tokens", args{"Hello, World!"}, 4},
		{"Count Chinese tokens", args{"你好，世界！"}, 7},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CountTokens(tt.args.text); got != tt.want {
				t.Errorf("CountTokens() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMakeFilename(t *testing.T) {
	type args struct {
		infile string
		ext    string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"Make summary file", args{"/tmp/myfile.txt", "sum"}, "/tmp/myfile.txt.sum"},
		{"Make json file", args{"/tmp/myfile.txt", "json"}, "/tmp/myfile.txt.json"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := MakeFilename(tt.args.infile, tt.args.ext); got != tt.want {
				t.Errorf("MakeFilename() = %v, want %v", got, tt.want)
			}
		})
	}
}
