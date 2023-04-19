package firebase_client

func NewBundledUpdate() *BundledUpdate {
	return &BundledUpdate{
		Ready:           false,
		InvocationCount: make(map[string]int64),
		Registrations:   make(map[string]RegistrationAction),
		Cache:           make(map[string]string),
	}
}
