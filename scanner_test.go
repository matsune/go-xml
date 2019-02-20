package xml

import (
	"testing"
)

func TestScanner_Get(t *testing.T) {
	type fields struct {
		source string
		cursor int
	}
	tests := []struct {
		name   string
		fields fields
		want   rune
	}{
		{
			name: "empty",
			fields: fields{
				source: "",
				cursor: 0,
			},
			want: 0,
		},
		{
			name: "get success",
			fields: fields{
				source: "a",
				cursor: 0,
			},
			want: 'a',
		},
		{
			name: "out of bounds",
			fields: fields{
				source: "a",
				cursor: 1,
			},
			want: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Scanner{
				source: []rune(tt.fields.source),
				cursor: tt.fields.cursor,
			}
			if got := s.Get(); got != tt.want {
				t.Errorf("Scanner.Get() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_Test(t *testing.T) {
	tests := []struct {
		name   string
		cursor int
		r      rune
		want   bool
	}{
		{
			cursor: 0,
			r:      ' ',
			want:   true,
		},
		{
			cursor: 1,
			r:      'a',
			want:   true,
		},
		{
			cursor: 2,
			r:      'あ',
			want:   true,
		},
		{
			cursor: 3,
			r:      'あ',
			want:   false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Scanner{
				source: []rune(" aあ"),
				cursor: tt.cursor,
			}
			if got := s.Test(tt.r); got != tt.want {
				t.Errorf("Scanner.Test() = %v, want %v", got, tt.want)
			}
		})
	}
}
