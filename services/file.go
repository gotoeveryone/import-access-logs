package services

import (
	"archive/zip"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
)

var (
	tempDir string
)

// CreateTempDir 一時ディレクトリを作成します。
func CreateTempDir(dirname string) error {
	tempDir = dirname + "tmp/"
	return os.MkdirAll(tempDir, 0775)
}

// GetFile ファイルを取得します。
func GetFile(dirname, filename string) (*os.File, error) {
	// ZIPファイルを開く
	r, err := zip.OpenReader(dirname + filename)
	if err != nil {
		return nil, err
	}
	defer r.Close()

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

		file, err := os.Create(fmt.Sprintf("%s%s.log", tempDir, filename))
		if err != nil {
			return nil, err
		}

		if err := ioutil.WriteFile(file.Name(), stream, 0775); err != nil {
			return nil, err
		}

		// 現在は1つなので1周で終了
		return file, nil
	}
	return nil, errors.New("対象ファイルなし")
}

// RemoveTempDir 一時ディレクトリを削除します。
func RemoveTempDir() error {
	if tempDir == "" {
		return nil
	}
	return os.RemoveAll(tempDir)
}
