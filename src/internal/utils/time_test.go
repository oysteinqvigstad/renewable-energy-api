package utils

import (
	"log"
	"testing"
)

func TestTime(t *testing.T) {
	ResetUptime()
	if GetUptime() > 1 {
		log.Fatal("ResetUptime did not reset the time. Was: ", GetUptime())
	}
}
