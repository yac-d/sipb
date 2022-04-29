package utils

import (
	"os"
	"log"
	"bytes"
	"io"
	"strings"
)

func FileExists(path string) bool {
	_, err := os.Stat(path)
	if err == nil { return true } 
	if os.IsNotExist(err) { return false } else { log.Fatalln(err) }
	return false
}

func ReaderLen(r io.Reader) int {
	var buf bytes.Buffer
	buf.ReadFrom(r)
	return len(buf.Bytes())
}

func StringsStartingWith(arr []string, substr string) []string {
	filtered := make([]string, 0)
	for _, str := range arr {
		if strings.HasPrefix(str, substr) {
			filtered = append(filtered, str)
		}
	}
	return filtered
}
