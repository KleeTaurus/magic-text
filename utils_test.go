package magictext

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

func BenchmarkCountTokens(b *testing.B) {
	text := "This is an era where AI breakthrough is coming daily. We didn’t have many AI-generated in public a few years ago, but now the technology is accessible to everyone. It’s excellent for many individual creators or companies that want to significantly take advantage of the technology to develop something complex, which might take a long time."
	for i := 0; i < b.N; i++ {
		CountTokens(text)
	}
}
