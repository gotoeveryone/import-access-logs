package services

import (
	"add-access-detail/models"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gocraft/dbr"
	"github.com/gotoeveryone/golang/common"
	"github.com/gotoeveryone/golang/logs"
)

// CreateSession コネクションを生成します。
func CreateSession(config common.Config) (*dbr.Tx, error) {
	// コネクションオープン
	connStr := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=true&loc=Local",
		config.DB.User, config.DB.Password, config.DB.Host,
		config.DB.Port, config.DB.Name)
	conn, _ := dbr.Open("mysql", connStr, nil)
	sess := conn.NewSession(nil)
	return sess.Begin()
}

// Exist データが存在するかを確認します。
func Exist(tx *dbr.Tx, detail models.AccessDetail) (bool, error) {
	cond := dbr.And(
		dbr.Eq("ip_address", detail.IPAddress),
		dbr.Eq("access_time", detail.AccessTime),
		dbr.Eq("access_url", detail.AccessURL),
		dbr.Eq("user_agent", detail.UserAgent))

	var cnt models.AccessDetail
	if err := tx.Select("*").From("access_logs").Where(cond).LoadStruct(&cnt); err != nil {
		if err != dbr.ErrNotFound {
			// 該当ありとみなす
			return true, err
		}
	}

	return (cnt.ID > 0), nil
}

// Insert データを登録し、処理成功した件数を返却します。
func Insert(tx *dbr.Tx, details *[]models.AccessDetail) (int, error) {
	// IDからデータを取得
	cnt := 0
	for _, detail := range *details {
		// 存在確認
		exist, err := Exist(tx, detail)
		if err != nil {
			return 0, err
		} else if exist {
			logs.Info(fmt.Sprintf("すでに存在します。【%s】【%s】【%s】【%s】",
				detail.IPAddress, detail.AccessTime, detail.AccessURL, detail.UserAgent))
			continue
		}

		// 登録
		if err := insert(tx, detail); err != nil {
			return cnt, err
		}
		cnt++
	}
	return cnt, nil
}

// データを保存します。
func insert(tx *dbr.Tx, detail models.AccessDetail) error {
	now := time.Now()
	detail.Created = &now

	_, err := tx.InsertInto("access_logs").
		Columns("id", "site_type", "ip_address", "access_date", "access_time", "access_url", "http_referer", "user_agent", "created").
		Record(detail).Exec()

	if err != nil {
		logs.Error("登録エラー")
		logs.Error(detail)
		return err
	}

	return nil
}
