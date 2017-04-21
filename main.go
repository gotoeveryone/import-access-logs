package main

import (
	"add-access-detail/controllers"
	"add-access-detail/models"
	"add-access-detail/services"
	"bytes"
	"fmt"
	"strings"
	"time"

	"os"

	"github.com/gotoeveryone/golang/common"
	"github.com/gotoeveryone/golang/logs"
	"github.com/gotoeveryone/golang/mail"
)

const (
	path       = "/share/analytics/k2ss.info/logs/"
	successDir = "/share/analytics/完了/"
)

var (
	// 設定
	config common.Config

	// 出力用
	results []models.Result
)

func main() {
	start := time.Now()

	// 設定ファイル読み出し
	common.LoadConfig(&config)

	// 件名
	subject := fmt.Sprintf("【自動通知】%s_K2SSアクセスデータ取り込み", time.Now().Format("20060102"))

	// エラー
	errors := []string{}

	// 一時ディレクトリ作成
	if err := services.CreateTempDir(path); err != nil {
		logs.Error(err)
		errors = append(errors, err.Error())
	} else {
		// アクセスログ登録
		if err := controllers.AddAccessLogs(config, path, successDir, &results); err != nil {
			logs.Error(err)
			errors = append(errors, err.Error())
			subject = "※失敗！ " + subject
		}

		// 一時ディレクトリ削除
		if err := services.RemoveTempDir(); err != nil {
			logs.Error(err)
			errors = append(errors, err.Error())
		}
	}

	// メール送信
	if len(results) > 0 || len(errors) > 0 {
		if err := sendMail(start, subject, errors); err != nil {
			logs.Error(err)
			os.Exit(1)
		}
	}

	// エラー保持状態なら異常終了
	if len(errors) > 0 {
		os.Exit(1)
	}
}

// 結果通知メールを送信します。
func sendMail(start time.Time, subject string, errors []string) error {
	var buffer bytes.Buffer
	buffer.WriteString(fmt.Sprintf("開始時間：%s | 終了時間：%s",
		start.Format("2006/01/02 15:04:05"), time.Now().Format("2006/01/02 15:04:05")))
	buffer.WriteString("\n\n")

	// エラーがあれば追加
	if len(errors) > 0 {
		buffer.WriteString(strings.Join(errors, "\n"))
		buffer.WriteString("\n\n")
	}

	// 本文作成
	for _, v := range results {
		buffer.WriteString(fmt.Sprintf("%s - %d件\n", v.Day, v.Count))
	}

	// メール送信
	if err := mail.SendMail(config, subject, buffer.String()); err != nil {
		return err
	}
	logs.Info("メールを送信しました。")
	return nil
}
