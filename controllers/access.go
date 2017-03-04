package controllers

import (
	"add-access-detail/models"
	"add-access-detail/services"
	"bufio"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/gocraft/dbr"
	"github.com/gotoeveryone/golang/common"
	"github.com/gotoeveryone/golang/logs"
)

var (
	pattern = regexp.MustCompile(".* ([0-9\\.]*) - - \\[(.*)\\] \"GET (.*) HTTP/[0-9\\.]{3}\" ([0-9]{3}) [0-9-]* \"(.*)\" \"(.*)\" .*")
	static  = regexp.MustCompile(".*(css|ico|js|img).*")
	day     = regexp.MustCompile(".*([0-9]{8}).*")
)

// AddAccessLogs アクセスログを登録します。
func AddAccessLogs(config common.Config, path, successDir string, results *[]models.Result) error {
	// ファイル一覧を取得
	files, err := ioutil.ReadDir(path)
	if err != nil {
		return err
	}

	// 取得したファイルの数だけ処理
	for _, fileInfo := range files {
		// 処理対象ファイルが取得できなければ次へ
		file := getTarget(path, fileInfo)
		if file == nil {
			continue
		}
		defer file.Close()

		// 対象から日付が抽出できなければ次へ
		target := day.FindStringSubmatch(file.Name())
		if len(target) == 0 {
			logs.Error(err)
			continue
		}
		targetDay := target[1]

		// トランザクションの開始
		tx, err := services.CreateSession(config)
		if err != nil {
			return err
		}

		// データ登録
		if err := regist(tx, targetDay, file, results); err != nil {
			tx.Rollback()
			logs.Error("トランザクションをロールバックしました。")
			return err
		}

		tx.Commit()
		logs.Info("トランザクションをコミットしました。")

		// ファイルの移動
		backupDir := fmt.Sprintf("%s%s-%s/", successDir, targetDay[:4], targetDay[4:6])
		if err := os.MkdirAll(backupDir, 0775); err != nil {
			return err
		}
		if err := os.Rename(path+fileInfo.Name(), backupDir+fileInfo.Name()); err != nil {
			return err
		}
	}
	return nil
}

// regist データを登録します。
func regist(tx *dbr.Tx, targetDay string, file *os.File, results *[]models.Result) error {

	// アクセスログ情報を生成
	details, err := createDetails(tx, file)
	if err != nil {
		return err
	}

	// データベースへ登録
	cnt, err := services.Insert(tx, &details)
	if err != nil {
		return err
	}

	// 結果を格納
	*results = append(*results, models.Result{targetDay, cnt})
	logs.Info(fmt.Sprintf("%sのデータ：%d件", targetDay, cnt))

	return nil
}

// 処理対象かどうかを判定します。
func getTarget(path string, fileInfo os.FileInfo) *os.File {
	// ディレクトリは除外
	if fileInfo.IsDir() {
		return nil
	}

	// 圧縮ファイルからファイルが抽出できなければ除外
	file, err := services.GetFile(path, fileInfo.Name())
	if err != nil {
		logs.Error(err)
		return nil
	}

	return file
}

// 登録データを生成します。
func createDetails(tx *dbr.Tx, file *os.File) ([]models.AccessDetail, error) {
	lines := []string{}
	sc := bufio.NewScanner(file)
	for sc.Scan() {
		lines = append(lines, sc.Text())
	}

	// ファイル内を走査
	details := []models.AccessDetail{}
	for _, l := range lines {
		// 指定パターンにマッチしなければ除外
		m := pattern.FindStringSubmatch(l)
		if len(m) == 0 {
			continue
		}
		// 静的ファイルは除外
		if res := static.FindStringSubmatch(m[3]); len(res) > 0 {
			continue
		}
		// "/wp"は除外
		if strings.HasPrefix(m[3], "/wp/") || strings.HasPrefix(m[3], "/wp-") ||
			strings.HasPrefix(m[3], "/adm") || strings.HasPrefix(m[3], "/blog") {
			continue
		}
		// ステータスコード200以外は除外
		if m[4] != strconv.Itoa(http.StatusOK) {
			continue
		}

		// アクセス日時をパース
		t, err := time.Parse("02/Jan/2006:15:04:05 -0700", m[2])
		if err != nil {
			return nil, err
		}

		// モデルに格納
		access := models.AccessDetail{}
		access.IPAddress = m[1]
		if strings.HasPrefix(m[3], "/igo/") {
			access.SiteType = "Go to Everyone!"
		} else {
			access.SiteType = "K2SS"
		}
		access.AccessDate = &t
		access.AccessTime = &t
		access.AccessURL = m[3]
		access.HTTPReferer = m[5]
		access.UserAgent = m[6]

		details = append(details, access)
	}

	return details, nil
}

// ファイルを取得します。
func getFile(filename string) (*os.File, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	return file, nil
}
