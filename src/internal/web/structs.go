package web

// APIStatus holds the status information for various API components, webhook count, version, and uptime.
type APIStatus struct {
	Countriesapi    int    `json:"countries_api"`
	Notification_db int    `json:"notification_db"`
	Webhooks        int    `json:"webhooks"`
	Version         string `json:"version"`
	Uptime          int    `json:"uptime"`
}
type WebhookResponse struct {
	WebhookID string `json:"webhook_id"`
	URL       string `json:"url,omitempty"`
	Country   string `json:"country"`
	Calls     int64  `json:"calls"`
}
