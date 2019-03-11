package utils

import (
	"encoding/hex"
	"time"
)

func EncodeTimeStr(day time.Time) string {
	return hex.EncodeToString([]byte(day.Format("2006-01-02 15:04:05")))
}
