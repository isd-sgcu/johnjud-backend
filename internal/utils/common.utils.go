package utils

import "strings"

func TrimInList(word string, sep string, trimList map[string]struct{}) string {
	splitWordList := strings.Split(word, sep)
	if _, ok := trimList[splitWordList[1]]; !ok {
		return word
	}

	return strings.TrimPrefix(word, sep+splitWordList[1])
}

func IsExisted(e map[string]struct{}, key string) bool {
	_, ok := e[key]
	if ok {
		return true
	}
	return false
}
