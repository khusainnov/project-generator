package main

import (
	"runtime"
	"strconv"
	"strings"
)

const fallbackGoVersion = "1.23"

func goVersion() string {
	return parseGoVersion(runtime.Version())
}

func parseGoVersion(raw string) string {
	fields := strings.Fields(raw)
	if len(fields) == 0 {
		return fallbackGoVersion
	}

	parts := strings.SplitN(strings.TrimPrefix(fields[0], "go"), ".", 3)
	if len(parts) < 2 {
		return fallbackGoVersion
	}

	for _, p := range parts[:2] {
		if _, err := strconv.Atoi(p); err != nil {
			return fallbackGoVersion
		}
	}

	return parts[0] + "." + parts[1]
}
