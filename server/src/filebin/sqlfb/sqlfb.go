package sqlfb

import (
	"database/sql"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path"
	"time"

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

	_, err = fb.db.Exec("CALL INSERT_FILE(?, ?, ?, ?, ?)", id, h.Filename, filepath, written, mimetype)
	if err != nil {
		result.Error = err
		return result
	}

	result.Error = fb.RemoveOldFiles()
	return result
}

func (fb *SQLFileBin) RemoveOldFiles() error {
	cntRes := fb.Count()
	if cntRes.Error != nil {
		return cntRes.Error
	}
	if cntRes.Count > fb.config.MaxFileCnt && fb.config.MaxFileCnt > 0 {
		details, res := fb.DetailsOfNthNewest(fb.config.MaxFileCnt)
		if res.Error != nil {
			return res.Error
		}

		rows, err := fb.db.Query(
			`SELECT item.id, fileitem.location
			FROM item INNER JOIN fileitem ON item.id = fileitem.ID
			WHERE item.ts < ?`,
			details.Timestamp,
		)
		if err != nil {
			return err
		}
		for rows.Next() {
			var uuid, loc string
			rows.Scan(&uuid, &loc)
			err = os.Remove(loc)
			if err != nil {
				return err
			}
			_, err = fb.db.Exec("DELETE FROM item WHERE id = ?", uuid)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (fb *SQLFileBin) Count() (result filebin.FileCountResult) {
	row := fb.db.QueryRow("SELECT COUNT(*) FROM item")
	result.Error = row.Scan(&result.Count)
	return
}

func (fb *SQLFileBin) DetailsOfNthNewest(n int) (fd filedetails.FileDetails, result filebin.FileDetailsResult) {
	row := fb.db.QueryRow("CALL NTH_MOST_RECENT_FILE(?)", n-1)

	var timestampStr string
	result.Error = row.Scan(&fd.ID, &fd.Name, &fd.Location, &fd.Size, &fd.Type, &timestampStr)
	if result.Error != nil {
		return
	}

	result.Filename = fd.Name
	fd.Timestamp, result.Error = time.Parse("2006-01-02 15:04:05", timestampStr)
	return
}
