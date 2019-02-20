package xml

import (
	"testing"
)

func Test_isNum(t *testing.T) {
	tests := []struct {
		name string
		r    rune
		want bool
	}{
		{
			r:    '0',
			want: true,
		},
		{
			r:    '1',
			want: true,
		},
		{
			r:    '2',
			want: true,
		},
		{
			r:    '3',
			want: true,
		},
		{
			r:    '4',
			want: true,
		},
		{
			r:    '5',
			want: true,
		},
		{
			r:    '6',
			want: true,
		},
		{
			r:    '7',
			want: true,
		},
		{
			r:    '8',
			want: true,
		},
		{
			r:    '9',
			want: true,
		},
		{
			r:    'a',
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isNum(tt.r); got != tt.want {
				t.Errorf("isNum() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_isLetter(t *testing.T) {
	tests := []struct {
		name string
		r    rune
		want bool
	}{
		{
			r:    'a',
			want: true,
		},
		{
			r:    'あ',
			want: true,
		},
		{
			r:    '(',
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isLetter(tt.r); got != tt.want {
				t.Errorf("isLetter() = %v, want %v", got, tt.want)
			}
		})
	}
}
