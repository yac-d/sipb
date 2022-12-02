package sqlfb

import (
	"database/sql"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"os"
	"path"

	"github.com/Eeshaan-rando/sipb/src/configdef"
	"github.com/Eeshaan-rando/sipb/src/filebin"
	"github.com/Eeshaan-rando/sipb/src/filebin/filedetails"
	"github.com/Eeshaan-rando/sipb/src/utils"
	_ "github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
)

type SQLFileBin struct {
	config configdef.Config
	db     *sql.DB
}

func New(c configdef.Config) *SQLFileBin {
	var fb = SQLFileBin{config: c}
	if !utils.FileExists(c.BinDir) {
		os.MkdirAll(c.BinDir, 0755)
	}
	return &fb
}

func (fb *SQLFileBin) Initialize() (err error) {
	fb.db, err = sql.Open(
		"mysql",
		fmt.Sprintf("%s:%s@tcp(%s)/sipb", fb.config.DBUser, fb.config.DBPasswd, fb.config.DBHost),
	)
	return err
}

func (fb *SQLFileBin) SaveFile(f multipart.File, h *multipart.FileHeader) filebin.SaveFileResult {
	id := uuid.New().String()
	filepath := path.Join(fb.config.BinDir, id)

	var result = filebin.SaveFileResult{
		TruncatedBytes: 0,
		Error:          nil,
		Filename:       h.Filename,
		Location:       filepath,
	}

	persistedFile, err := os.Create(filepath)
	if err != nil {
		result.Error = err
		return result
	}

	var written int64
	if fb.config.MaxFileSize > -1 {
		written, err = io.CopyN(persistedFile, f, fb.config.MaxFileSize)
		realFileSize := int64(utils.ReaderLen(f))
		if realFileSize > fb.config.MaxFileSize {
			result.TruncatedBytes = realFileSize - fb.config.MaxFileSize
		}
	} else {
		written, err = io.Copy(persistedFile, f)
	}
	persistedFile.Close()

	mimetype, err := utils.MimeType(filepath)
	if err != nil {
		result.Error = err
		return result
	}

	log.Println(id, h.Filename, filepath, written, mimetype)
	_, err = fb.db.Exec("CALL INSERT_FILE(?, ?, ?, ?, ?)", id, h.Filename, filepath, written, mimetype)
	if err != nil {
		result.Error = err
		return result
	}

	result.Error = fb.RemoveOldFiles()

	return result
}

func (fb *SQLFileBin) RemoveOldFiles() error {
	return nil
}

func (fb *SQLFileBin) Count() (result filebin.FileCountResult) {
	row := fb.db.QueryRow("SELECT COUNT(*) FROM item")
	result.Error = row.Scan(&result.Count)
	return
}

func (fb *SQLFileBin) DetailsOfNthNewest(n int) (fd filedetails.FileDetails, result filebin.FileDetailsResult) {
	return
}
