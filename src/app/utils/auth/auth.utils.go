package auth

import (
	"fmt"
)

func FormatPath(method string, path string) string {
	return fmt.Sprintf("%v %v", method, path)
}
