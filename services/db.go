package services

import (
	"add-access-detail/models"
	"strconv"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gocraft/dbr"
	"github.com/gotoeveryone/golang/common"
)

// Exist データが存在するかを確認します。
func Exist(config common.Config, detail models.AccessDetail) (int, error) {
	ses := createSession(config)
	cond := dbr.And(
		dbr.Eq("ip_address", detail.IPAddress),
		dbr.Eq("access_time", detail.AccessTime),
		dbr.Eq("access_url", detail.AccessURL),
		dbr.Eq("user_agent", detail.UserAgent))

	var cnt int
	if err := ses.Select("id").From("access_logs").Where(cond).LoadValue(&cnt); err != nil && err != dbr.ErrNotFound {
		return 0, err
	}

	return cnt, nil
}

// Regist データを登録し、処理成功した件数を返却します。
func Regist(config common.Config, details *[]models.AccessDetail) (int, error) {
	ses := createSession(config)
	tx, _ := ses.Begin()

	// IDからデータを取得
	cnt := 0
	for _, detail := range *details {
		if err := save(tx, detail); err != nil {
			tx.Rollback()
			return cnt, err
		}
		cnt++
	}
	tx.Commit()

	return cnt, nil
}

// コネクションを生成します。
func createSession(config common.Config) *dbr.Session {
	// コネクションオープン
	connStr := config.DB.User + ":" + config.DB.Password + "@tcp(" + config.DB.Host + ":" + strconv.Itoa(config.DB.Port) + ")/" + config.DB.Name
	conn, _ := dbr.Open("mysql", connStr, nil)
	return conn.NewSession(nil)
}

// データを保存します。
func save(tx *dbr.Tx, detail models.AccessDetail) error {
	detail.Created = time.Now()

	_, err := tx.InsertInto("access_logs").
		Columns("id", "ip_address", "access_time", "access_url", "http_referer", "user_agent", "created").
		Record(detail).Exec()

	if err != nil {
		return err
	}

	return nil
}
