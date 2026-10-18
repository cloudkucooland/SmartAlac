package sa

import (
	"testing"
)

func TestSanitize(t *testing.T) {
	cases := []struct {
		in, want string
	}{
		{"Normal Title", "Normal Title"},
		{"AC/DC", "AC_DC"},
		{"Is This Love?", "Is This Love_"},
		{"<Angle> | Brackets", "_Angle_ _ Brackets"},
		{"File:Name", "File_Name"},
		{"Quote's and \"Quotes\"", "Quote_s and _Quotes_"},
		{"Vertical|Bar", "Vertical_Bar"},
		{"Asterisk*", "Asterisk_"},
	}

	for _, c := range cases {
		got := sanitize(c.in)
		if got != c.want {
			t.Errorf("sanitize(%q) == %q, want %q", c.in, got, c.want)
		}
	}
}
