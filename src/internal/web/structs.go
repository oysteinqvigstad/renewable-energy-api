package web

// APIStatus holds the status information for various API components, webhook count, version, and uptime.
type APIStatus struct {
	Countriesapi    int    `json:"countriesapi"`
	Notification_db int    `json:"notification_db"`
	Webhooks        int    `json:"webhooks"`
	Version         string `json:"version"`
	Uptime          int    `json:"uptime"`
}
