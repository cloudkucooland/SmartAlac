package sa

import "testing"

func TestMediumFormat(t *testing.T) {
	cases := []struct {
		in, want string
	}{
		{"12\" Vinyl", "12″ Vinyl"},
		{"7\" Vinyl", "7″ Vinyl"},
		{"10\" Vinyl", "10″ Vinyl"},
		{"CD", "CD"},
		{"Digital Media", "Digital Media"},
		{"5\" CD", "5″ CD"},
	}

	for _, c := range cases {
		got := mediumFormat(c.in)
		if got != c.want {
			t.Errorf("mediumFormat(%q) == %q, want %q", c.in, got, c.want)
		}
	}
}
