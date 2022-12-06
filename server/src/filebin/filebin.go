package filebin

import (
	"mime/multipart"

	"github.com/yac-d/sipb/filedetails"
)

type FileBin interface {
	SaveFile(toSave FileToSave) SaveFileResult
	DetailsOfNthNewest(n int) (filedetails.FileDetails, FileDetailsResult)
	Count() FileCountResult
	RemoveOldFiles() error
}

// FileToSave holds a file and its related information
// to be passed to the file bin for saving
type FileToSave struct {
	File   multipart.File
	Header *multipart.FileHeader
	Note   string
}

type SaveFileResult struct {
	Error          error
	TruncatedBytes int64
	Filename       string
	Location       string
}

type FileDetailsResult struct {
	Error    error
	Filename string
}

type FileCountResult struct {
	Error error
	Count int
}
