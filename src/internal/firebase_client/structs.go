package firebase_client

type InvocationRegistration struct {
	WebhookID string `json:"webhook_id"`
	URL       string `json:"url"`
	Country   string `json:"country"`
	Calls     int64  `json:"calls"`
}

type Subscriptions struct {
	invocationCount map[string]int64
	registrations   map[string]InvocationRegistration
}

type BundledUpdate struct {
	Ready           bool
	InvocationCount map[string]int64
	Registrations   map[string]RegistrationAction
	Cache           map[string]string
}

type RegistrationAction struct {
	Add          bool
	Registration InvocationRegistration
}