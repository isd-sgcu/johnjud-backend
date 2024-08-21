package auth

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/google/uuid"
)

func IsExisted(e map[string]struct{}, key string) bool {
	_, ok := e[key]
	return ok
}

func FormatPath(method string, path string, keys []string) string {
	for _, key := range keys {
		path = strings.Replace(path, key, ":id", 1)
	}

	return fmt.Sprintf("%v %v", method, path)
}

func FindIntFromStr(s string, sep string) []string {
	spliteds := strings.Split(s, sep)

	var result []string

	for _, splited := range spliteds {
		_, err := strconv.Atoi(splited)
		if err == nil {
			result = append(result, splited)
		}
	}

	return result
}

func FindUUIDFromStr(s string, sep string) []string {
	spliteds := strings.Split(s, sep)

	var result []string

	for _, splited := range spliteds {
		_, err := uuid.Parse(splited)
		if err == nil {
			result = append(result, splited)
		}
	}

	return result
}

func MergeStringSlice(s1 []string, s2 []string) []string {
	return append(s1, s2...)
}

func FindIDFromPath(path string) []string {
	uuids := FindUUIDFromStr(path, "/")
	ids := FindIntFromStr(path, "/")

	return MergeStringSlice(ids, uuids)
}
