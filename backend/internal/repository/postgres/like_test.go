package postgres

import "testing"

func TestEscapeLikePattern(t *testing.T) {
	cases := []struct {
		in   string
		want string
	}{
		{"atomic habits", "%atomic habits%"},
		{"100% original", `%100\% original%`},
		{"book_v2", `%book\_v2%`},
		{`path\to%_file`, `%path\\to\%\_file%`},
	}

	for _, tc := range cases {
		if got := escapeLikePattern(tc.in); got != tc.want {
			t.Errorf("escapeLikePattern(%q) = %q, want %q", tc.in, got, tc.want)
		}
	}
}
