package xml

import (
	"reflect"
	"testing"
)

func TestScanner_Get(t *testing.T) {
	type fields struct {
		source string
		cursor uint
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
			s := &scanner{
				source: []rune(tt.fields.source),
				cursor: tt.fields.cursor,
			}
			if got := s.Get(); got != tt.want {
				t.Errorf("Scanner.Get() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_scanner_Pos(t *testing.T) {
	type fields struct {
		source []rune
		cursor uint
	}
	tests := []struct {
		name   string
		fields fields
		want   Pos
	}{
		{
			fields: fields{
				source: []rune(`abc
def
ghijk
l`), // a b c \n d e f \n g h i j k \n l
				cursor: 10, // i
			},
			want: Pos{
				Line: 3,
				Col:  3,
			},
		},
		{
			fields: fields{
				source: []rune(`abc`),
				cursor: 2,
			},
			want: Pos{
				Line: 1,
				Col:  3,
			},
		},
		{
			name: "EOF",
			fields: fields{
				source: []rune(`abc`),
				cursor: 3,
			},
			want: Pos{
				Line: 1,
				Col:  4,
			},
		},
		{
			name: "out of bounds",
			fields: fields{
				source: []rune(`abc`),
				cursor: 4,
			},
			want: Pos{
				Line: 1,
				Col:  4,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &scanner{
				source: tt.fields.source,
				cursor: tt.fields.cursor,
			}
			if got := s.pos(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("scanner.Position() = %v, want %v", got, tt.want)
			}
		})
	}
}
