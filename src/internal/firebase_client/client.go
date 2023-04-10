package firebase_client

import (
	"assignment2/internal/datastore"
	"time"
)

type InvocationRegistration struct {
	WebhookID string `json:"webhook_id"`
	URL       string `json:"url"`
	Country   string `json:"country"`
	Calls     int    `json:"calls"`
}

type FirebaseClient struct {
}

func NewFirebaseClient() (*FirebaseClient, error) {
	return &FirebaseClient{}, nil
}

func (client *FirebaseClient) SetInvocationCount(ccna3 string, number int) {
	// e.g. SetInvocationCount("NOR", 5)
}

func (client *FirebaseClient) GetInvocationCount(ccna3 string) {
	// e.g. GetInvocationCount("NOR")
}

func (client *FirebaseClient) SetRenewablesCache(url string, list *datastore.YearRecordList) {
	// e.g. SetRenewablesCache("/current/nor?neighbours=true", *data)
}

func (client *FirebaseClient) GetRenewablesCache(url string) (datastore.YearRecordList, time.Time, error) {
	// e.g. SetRenewablesCache("/current/nor?neighbours=true")
	return datastore.YearRecordList{}, time.Time{}, nil
}

func (client *FirebaseClient) DeleteRenewablesCache(url string) {
	// e.g. DeleteRenewablesCache("/current/nor?neighbours=true")
}

func (client *FirebaseClient) SetInvocationRegistration(registration InvocationRegistration) {
	// e.g. SetInvocationRegistration(registration)
}

func (client *FirebaseClient) GetAllInvocationRegistrations() []InvocationRegistration {
	// e.g. GetAllInvocationRegistrations()
	return []InvocationRegistration{}
}
