package types

// InvocationRegistration represents a webhook registration with its associated information.
type InvocationRegistration struct {
	WebhookID string `json:"webhook_id"`
	URL       string `json:"url"`
	Country   string `json:"country"`
	Calls     int64  `json:"calls"`
}

// RegistrationAction represents an action to add or remove a webhook registration.
type RegistrationAction struct {
	Add          bool
	Registration InvocationRegistration
}

// BundledUpdate represents a set of updates to be performed, including invocation counts, registrations, and cache updates.
type BundledUpdate struct {
	Ready           bool
	InvocationCount map[string]int64
	Registrations   map[string]RegistrationAction
	Cache           map[string]YearRecordList
}
