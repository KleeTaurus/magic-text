package summaryit

import "testing"

func TestCountTokens(t *testing.T) {
	type args struct {
		text string
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{"English Tokens", args{"Hello, World!"}, 4},
		{"Chinese Tokens", args{"你好，世界！"}, 7},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CountTokens(tt.args.text); got != tt.want {
				t.Errorf("CountTokens() = %v, want %v", got, tt.want)
			}
		})
	}
}
