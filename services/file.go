package services

import (
	"archive/zip"
	"errors"
	"io/ioutil"
	"os"
)

// GetFile ファイルを取得します。
func GetFile(filename string) (*os.File, error) {
	// ZIPファイルを開く
	r, err := zip.OpenReader("/share/analytics/k2ss.info/logs/ssl_access_log." + filename + ".zip")
	if err != nil {
		return nil, err
	}
	// 内部のファイルを走査
	for _, f := range r.File {
		// ファイルを開く
		rc, err := f.Open()
		if err != nil {
			return nil, err
		}
		defer rc.Close()

		stream, err := ioutil.ReadAll(rc)
		if err != nil {
			return nil, err
		}

		if err := os.MkdirAll("/share/analytics/backup/", 0775); err != nil {
			return nil, err
		}

		file, err := os.Create("/share/analytics/backup/" + filename + ".log")
		if err != nil {
			return nil, err
		}

		if err := ioutil.WriteFile(file.Name(), stream, 0775); err != nil {
			return nil, err
		}

		// 現在は1つなので1周で終了
		return file, nil
	}
	return nil, errors.New("走査対象なし")
}
