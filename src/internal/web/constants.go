package web

const (
	Version               = "v1"
	DefaultPath           = "/energy/" + Version + "/"
	RenewablesCurrentPath = DefaultPath + "renewables/current/"
	RenewablesHistoryPath = DefaultPath + "renewables/history/"
	NotificationsPath     = DefaultPath + "notifications/"
	StatusPath            = DefaultPath + "status/"
	FirebaseUpdateFreq    = 5 // update firebase every 5 seconds
)
