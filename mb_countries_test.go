package sa

import "testing"

func TestResolveCountry(t *testing.T) {
	cases := []struct {
		in, want string
	}{
		{"US", "United States"},
		{"GB", "United Kingdom"},
		{"DE", "Germany"},
		{"JP", "Japan"},
		{"XX", ""}, // Unknown
		{"", ""},   // Empty
	}

	for _, c := range cases {
		got := resolveCountry(c.in)
		if got != c.want {
			t.Errorf("resolveCountry(%q) == %q, want %q", c.in, got, c.want)
		}
	}
}
