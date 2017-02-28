package main

import (
	"add-access-detail/controllers"
	"bytes"
	"strconv"
	"time"

	"github.com/gotoeveryone/golang/common"
	"github.com/gotoeveryone/golang/logs"
	"github.com/gotoeveryone/golang/mail"
)

func main() {
	// 設定ファイル読み出し
	var config common.Config
	common.LoadConfig(&config)

	// 件名
	subject := "【自動通知】" + time.Now().Format("20060102") + "_K2SSアクセスデータ取り込み"

	// データ登録
	results := map[string]int{}
	if err := controllers.Regist(config, &results); err != nil {
		logs.Error(err)
		subject = subject + "異常終了"
	}

	// 本文
	var buffer bytes.Buffer
	for k, v := range results {
		buffer.WriteString(k + ": " + strconv.Itoa(v) + "件\n")
	}

	// メール送信
	if err := mail.SendMail(config, subject, buffer.String()); err != nil {
		logs.Error(err)
	}
	logs.Info("メールを送信しました。")
}
