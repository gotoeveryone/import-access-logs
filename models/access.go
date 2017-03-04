package models

import "time"

type (
	// AccessDetail アクセス詳細
	AccessDetail struct {
		ID          int        `db:"id"`
		SiteType    string     `db:"site_type"`
		IPAddress   string     `db:"ip_address"`
		AccessDate  *time.Time `db:"access_date"`
		AccessTime  *time.Time `db:"access_time"`
		AccessURL   string     `db:"access_url"`
		HTTPReferer string     `db:"http_referer"`
		UserAgent   string     `db:"user_agent"`
		Created     *time.Time `db:"created"`
	}

	// Result 登録結果
	Result struct {
		Day   string
		Count int
	}
)
