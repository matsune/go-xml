package xml

import "testing"

func TestXMLDecl_String(t *testing.T) {
	type fields struct {
		Version    string
		Encoding   string
		Standalone bool
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			fields: fields{
				Version:    "1.0",
				Encoding:   "UTF-8",
				Standalone: true,
			},
			want: `<?xml version="1.0" encoding="UTF-8" standalone="yes" ?>`,
		},
		{
			fields: fields{
				Version: "1.1",
			},
			want: `<?xml version="1.1" standalone="no" ?>`,
		},
		{
			fields: fields{},
			want:   `<?xml version="" standalone="no" ?>`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			x := &XMLDecl{
				Version:    tt.fields.Version,
				Encoding:   tt.fields.Encoding,
				Standalone: tt.fields.Standalone,
			}
			if got := x.String(); got != tt.want {
				t.Errorf("XMLDecl.String() = %v, want %v", got, tt.want)
			}
		})
	}
}
