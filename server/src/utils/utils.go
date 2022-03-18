package utils

import (
	"os"
	"log"
)

func FileExists(path string) bool {
	_, err := os.Stat(path)
	if err == nil { return true } 
	if os.IsNotExist(err) { return false } else { log.Fatalln(err) }
	return false
}
