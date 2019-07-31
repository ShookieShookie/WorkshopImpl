package main

import "testing"

func Test_getCard(t *testing.T) {
	tests := []struct {
		name  string
		write string
		want  string
	}{
		{
			name:  "read from stdin",
			write: "hello",
			want:  "hello",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getCard(); got != tt.want {
				t.Errorf("getCard() = %v, want %v", got, tt.want)
			}
		})
	}
}
