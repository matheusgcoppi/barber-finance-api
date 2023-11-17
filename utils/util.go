package utils

import (
	"strings"
)

const (
	EmailPattern = `^([a-z\d\.-]+)@([a-z\d-]+)\.([a-z]{2,10})(\.[a-z]{2,8})?$`
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
