package utils

import (
	"bytes"
	"io"
	"log"
	"net/http"
	"os"
)

func FileExists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	} else {
		log.Fatalln(err)
	}
	return false
}

func ReaderLen(r io.Reader) int {
	var buf bytes.Buffer
	buf.ReadFrom(r)
	return len(buf.Bytes())
}

func MimeType(filepath string) (string, error) {
	f, err := os.Open(filepath)
	if err != nil {
		return "", err
	}
	defer f.Close()

	var fHeader = make([]byte, 512)
	f.Read(fHeader)

	return http.DetectContentType(fHeader), nil
}
