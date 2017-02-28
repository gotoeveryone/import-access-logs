package controllers

import (
	"add-access-detail/models"
	"add-access-detail/services"
	"bufio"
	"fmt"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"time"

	"github.com/gotoeveryone/golang/common"
	"github.com/gotoeveryone/golang/logs"
)

var (
	pattern = regexp.MustCompile(".* ([0-9\\.]*) - - \\[(.*)\\] \"GET (.*) HTTP/[0-9\\.]{3}\" ([0-9]{3}) [0-9-]* \"(.*)\" \"(.*)\" .*")
	sp      = regexp.MustCompile(".*[css|ico|js|img].*")
)

// Regist データを登録します。
func Regist(config common.Config, results *map[string]int) error {
	// 1月1日
	now := time.Now()
	jan := time.Date(now.Year(), time.January, 1, 0, 0, 0, 0, time.Local)
	feb := time.Date(now.Year(), time.February, 1, 0, 0, 0, 0, time.Local)
	for t := 0; t < feb.AddDate(0, 0, -1).Day(); t++ {
		target := jan.AddDate(0, 0, t)
		f := target.Format("20060102")
		logs.Info(f + "のデータ")
		file, err := services.GetFile(f)
		if err != nil {
			logs.Error(err)
			continue
		}
		details, err := createDetails(config, file)
		if err != nil {
			logs.Error(err)
			continue
		}
		if _, err := services.Regist(config, &details); err != nil {
			return err
		}
		(*results)[f] = len(details)

		// ファイル移動
		backupDir := "/share/analytics/完了/" + target.Format("2006-01") + "/"
		if err := os.MkdirAll(backupDir, 0775); err != nil {
			return err
		}
		if err := os.Rename("/share/analytics/k2ss.info/logs/ssl_access_log."+f+".zip", backupDir+"ssl_access_log."+f+".zip"); err != nil {
			return err
		}
	}
	return nil
}

// CreateDetails 登録データを生成します。
func createDetails(config common.Config, file *os.File) ([]models.AccessDetail, error) {
	lines := []string{}
	sc := bufio.NewScanner(file)
	for sc.Scan() {
		lines = append(lines, sc.Text())
	}

	// ファイル内を走査
	details := []models.AccessDetail{}
	for _, l := range lines {
		// log.Println(l)
		m := pattern.FindStringSubmatch(l)
		if len(m) == 0 {
			continue
		}
		// 静的ファイルは除外
		if res := sp.FindStringSubmatch(m[3]); len(res) > 0 {
			continue
		}
		// ステータスコード200以外は除外
		if m[4] != strconv.Itoa(http.StatusOK) {
			continue
		}
		t, err := time.Parse("02/Jan/2006:15:04:05 -0700", m[2])
		if err != nil {
			return nil, err
		}
		access := models.AccessDetail{}
		access.IPAddress = m[1]
		access.AccessTime = t
		access.AccessURL = m[3]
		access.HTTPReferer = m[5]
		access.UserAgent = m[6]

		cnt, err := services.Exist(config, access)
		if err != nil {
			return nil, err
		} else if cnt > 0 {
			logs.Info(fmt.Sprintf("すでに存在します。【%s】【%s】【%s】【%s】",
				access.IPAddress, access.AccessTime, access.AccessURL, access.UserAgent))
			continue
		}

		details = append(details, access)
	}
	logs.Info(strconv.Itoa(len(details)) + "件")
	return details, nil
}
