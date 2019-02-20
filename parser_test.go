package xml

import (
	"reflect"
	"testing"
)

func TestParser_Test(t *testing.T) {
	type fields struct {
		source string
		cursor int
	}
	tests := []struct {
		name   string
		fields fields
		str    string
		want   bool
	}{
		{
			fields: fields{
				source: "abc",
				cursor: 0,
			},
			str:  "a",
			want: true,
		},
		{
			fields: fields{
				source: "abc",
				cursor: 0,
			},
			str:  "ab",
			want: true,
		},
		{
			fields: fields{
				source: "abc",
				cursor: 0,
			},
			str:  "abc",
			want: true,
		},
		{
			fields: fields{
				source: "abc",
				cursor: 1,
			},
			str:  "a",
			want: false,
		},
		{
			fields: fields{
				source: "abc",
				cursor: 1,
			},
			str:  "bc",
			want: true,
		},
		{
			fields: fields{
				source: "abc",
				cursor: 2,
			},
			str:  "c",
			want: true,
		},
		{
			fields: fields{
				source: "abc",
				cursor: 2,
			},
			str:  "a",
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Parser{
				source: []rune(tt.fields.source),
				cursor: tt.fields.cursor,
			}
			if got := p.Tests(tt.str); got != tt.want {
				t.Errorf("Parser.Match() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParser_Get(t *testing.T) {
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
			fields: fields{
				source: "a",
				cursor: 0,
			},
			want: 'a',
		},
		{
			fields: fields{
				source: "a",
				cursor: 1,
			},
			want: EOF,
		},
		{
			fields: fields{
				source: "a",
				cursor: 2,
			},
			want: EOF,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Parser{
				source: []rune(tt.fields.source),
				cursor: tt.fields.cursor,
			}
			if got := p.Get(); got != tt.want {
				t.Errorf("Parser.Get() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParser_parseXmlDecl(t *testing.T) {

	tests := []struct {
		name    string
		str     string
		want    *XMLDecl
		wantErr bool
	}{
		{
			str: `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>`,
			want: &XMLDecl{
				VersionInfo:  "1.0",
				EncodingDecl: "UTF-8",
				Standalone:   true,
			},
		},
		{
			str: `<?xml version="2.0"  standalone="no"    ?>`,
			want: &XMLDecl{
				VersionInfo:  "2.0",
				EncodingDecl: "",
				Standalone:   false,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := NewParser(tt.str)
			got, err := p.parseXmlDecl()
			if (err != nil) != tt.wantErr {
				t.Errorf("Parser.parseXmlDecl() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Parser.parseXmlDecl() = %v, want %v", got, tt.want)
			}
		})
	}
}
