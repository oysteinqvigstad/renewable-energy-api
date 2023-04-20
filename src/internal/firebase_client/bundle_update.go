package firebase_client

import "assignment2/internal/datastore"

func NewBundledUpdate() *BundledUpdate {
	return &BundledUpdate{
		Ready:           false,
		InvocationCount: make(map[string]int64),
		Registrations:   make(map[string]RegistrationAction),
		Cache:           make(map[string]datastore.YearRecordList),
	}
}
