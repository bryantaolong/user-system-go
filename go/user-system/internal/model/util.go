package model

import "strings"

func splitByComma(s string) []string {
	return strings.Split(s, ",")
}

func trimSpace(s string) string {
	return strings.TrimSpace(s)
}
