package firebase_client

import (
	"assignment2/internal/types"
)

func NewBundledUpdate() *types.BundledUpdate {
	return &types.BundledUpdate{
		Ready:           false,
		InvocationCount: make(map[string]int64),
		Registrations:   make(map[string]types.RegistrationAction),
		Cache:           make(map[string]types.YearRecordList),
	}
}
