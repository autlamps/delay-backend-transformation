package main

import (
	"testing"
	"time"
)

// TODO: move this testing into a package which makes more sense
func TestParseTime(t *testing.T) {
	tests := []struct {
		in  string
		out time.Time
	}{
		{"13:45:01", time.Date(0, 1, 1, 13, 45, 01, 0, time.UTC)},
		{"24:00:01", time.Date(0, 1, 1, 00, 00, 01, 0, time.UTC)},
		{"26:35:19", time.Date(0, 1, 1, 02, 35, 19, 0, time.UTC)},
	}

	for _, test := range tests {
		parsed := ParseTime(test.in)

		if !test.out.Equal(parsed) {
			t.Errorf("%v - Expected %v, got %v", test.in, test.out, parsed)
		}
	}
}