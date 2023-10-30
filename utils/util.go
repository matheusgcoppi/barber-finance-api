package utils

import (
	"strings"
)

func ContainsUtil(words []string, json string) bool {
	length := len(words)
	count := 0
	for i := 0; i < length; i++ {
		if strings.Contains(json, words[i]) {
			count++
		}
	}
	if length == count {
		return true
	}
	return false
}
