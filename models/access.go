package models

import "time"

type (
	// AccessDetail アクセス詳細
	AccessDetail struct {
		ID          int       `db:"id"`
		IPAddress   string    `db:"ip_address"`
		AccessTime  time.Time `db:"access_time"`
		AccessURL   string    `db:"access_url"`
		HTTPReferer string    `db:"http_referer"`
		UserAgent   string    `db:"user_agent"`
		Created     time.Time `db:"created"`
	}
)
