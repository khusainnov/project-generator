package main

import "testing"

func TestParseGoVersion(t *testing.T) {
	tests := map[string]string{
		"go1.27.0":                       "1.27",
		"go1.27":                         "1.27",
		"go1.30.11":                      "1.30",
		"go1.21.0 X:nocoverageredesign":  "1.21",
		"go1.27rc1":                      fallbackGoVersion,
		"go1.28beta1":                    fallbackGoVersion,
		"devel go1.28-abc123 2026-01-01": fallbackGoVersion,
		"go1":                            fallbackGoVersion,
		"":                               fallbackGoVersion,
	}

	for raw, want := range tests {
		if got := parseGoVersion(raw); got != want {
			t.Errorf("parseGoVersion(%q) = %q, want %q", raw, got, want)
		}
	}
}
