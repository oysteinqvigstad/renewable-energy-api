package app

import "time"

var startTime time.Time

func ResetUptime() {
	startTime = time.Now()
}

func GetUptime() int {
	return int(time.Since(startTime).Seconds())
}
