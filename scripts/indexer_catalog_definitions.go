package main

import (
	"encoding/json"
	"os"
	"strings"
)

func writeDefinitions(path string, definitions map[string]string) error {
	data, err := json.MarshalIndent(definitions, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, append(data, '\n'), 0o644)
}

func cleanProwlarrYAML(body []byte) []byte {
	out := make([]byte, 0, len(body))
	inDouble := false
	for i := 0; i < len(body); i++ {
		c := body[i]
		if c == '"' && !oddTrailingBackslashes(out) {
			inDouble = !inDouble
			out = append(out, c)
			continue
		}
		if inDouble && c == '\\' && i+1 < len(body) {
			next := body[i+1]
			if next == '\\' {
				out = append(out, c, next)
				i++
				continue
			}
			if !isYAMLEscape(next) {
				out = append(out, '\\')
			}
		}
		out = append(out, c)
	}
	return out
}

func isYAMLEscape(c byte) bool {
	return strings.ContainsRune(`0abtnvfre "N_LPxuU`, rune(c))
}

func oddTrailingBackslashes(value []byte) bool {
	count := 0
	for i := len(value) - 1; i >= 0 && value[i] == '\\'; i-- {
		count++
	}
	return count%2 == 1
}
