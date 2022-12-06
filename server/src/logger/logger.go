package logger

import (
	"log"

	"github.com/yac-d/sipb/filebin"
)

func LogConfigRead(source string, err error) {
	if err != nil {
		log.Println("Error reading configuration from", source)
		log.Fatalln(err)
	}
	log.Println("Read configuration from", source)
}

func LogFileSave(result filebin.SaveFileResult) {
	log.Printf("Saving file %s", result.Filename)
	if result.Error != nil {
		log.Printf("Error while saving file: %s", result.Error)
	}
	if result.TruncatedBytes > 0 {
		log.Printf("File %s above set size limit, truncating by %d bytes", result.Filename, result.TruncatedBytes)
	}
}

func LogFileDetailsRequest(result filebin.FileDetailsResult) {
	log.Printf("Details for file %s requested", result.Filename)
	if result.Error != nil {
		log.Printf("Error while getting file details: %s")
	}
}

func LogFileCountRequest(result filebin.FileCountResult) {
	if result.Error != nil {
		log.Printf("Error while getting file count: %s", result.Error)
	} else {
		log.Printf("File count requested; currently at %d", result.Count)
	}
}

func Log(thing any) {
	log.Println(thing)
}
