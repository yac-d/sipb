package utils

import (
	"os"
	"log"
	"bytes"
	"io"
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
