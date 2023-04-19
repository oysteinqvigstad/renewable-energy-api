package web

const (
	DefaultPath           = "/energy/v1/"
	RenewablesCurrentPath = DefaultPath + "renewables/current/"
	RenewablesHistoryPath = DefaultPath + "renewables/history/"
	NotificationsPath     = DefaultPath + "notifications/"
	StatusPath            = DefaultPath + "status/"
	FirebaseUpdateFreq    = 5 // update firebase every 5 seconds
)
