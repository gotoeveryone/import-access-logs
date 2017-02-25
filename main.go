package main

import (
	"add-access-detail/models"
	"add-access-detail/services"
	"bufio"
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

func main() {
	// 設定ファイル読み出し
	var config common.Config
	common.LoadConfig(&config)

	// 1月1日
	now := time.Now()
	jan := time.Date(now.Year(), time.January, 1, 0, 0, 0, 0, time.Local)
	feb := time.Date(now.Year(), time.February, 1, 0, 0, 0, 0, time.Local)
	for t := 0; t < feb.AddDate(0, 0, -1).Day(); t++ {
		f := jan.AddDate(0, 0, t).Format("20060102")
		logs.Info(f + "のデータ")
		file, err := services.GetFile(f)
		if err != nil {
			logs.Error(err)
			continue
		}
		details, err := loadModel(file)
		if err != nil {
			logs.Error(err)
			continue
		}
		if _, err := services.Regist(config, &details); err != nil {
			logs.Fatal(err)
		}
	}
}

func loadModel(file *os.File) ([]models.AccessDetail, error) {
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
		access.AccessDate = t
		access.AccessURL = m[3]
		access.HTTPReferer = m[5]
		access.UserAgent = m[6]

		details = append(details, access)
	}
	logs.Info(strconv.Itoa(len(details)) + "件")
	return details, nil
}
