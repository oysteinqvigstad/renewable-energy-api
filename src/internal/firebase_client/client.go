package firebase_client

import (
	"assignment2/internal/datastore"
	"cloud.google.com/go/firestore" // Firestore-specific support
	"context"                       // State handling across API boundaries; part of native GoLang API
	firebase "firebase.google.com/go"
	"fmt"
	"google.golang.org/api/option"
	"log"
	"time"
)

/*
FirebaseClient is a wrapper around Firestore client that includes
the context necessary for performing operations on the Firestore.
*/
type FirebaseClient struct {
	// Firebase context and client used by Firestore functions throughout the program.
	ctx    context.Context
	client *firestore.Client
}

// NewFirebaseClient initializes and returns a new FirebaseClient by connecting to Firebase using the secret key.
func NewFirebaseClient() (*FirebaseClient, error) {
	/* Firebase initialisation  -> means setting up the connection to Firebase */
	ctx := context.Background()                                    // Create a basic empty box for tasks
	secret_key_path := option.WithCredentialsFile(PathToSecretKey) // Tell the program where to find the secret key for Firebase
	app, err := firebase.NewApp(ctx, nil, secret_key_path)         // Connect to Firebase using the secret key

	if err != nil { // If there's an error, stop the program and show the error
		log.Fatalf("Error initializing Firebase app: %v", err)
		return nil, err
	}

	/* Instantiate client */
	client, err := app.Firestore(ctx) // Get the helper for talking to Firestore
	if err != nil {
		log.Fatalf("Failed to create Firestore client: %v", err)
		return nil, err
	}

	return &FirebaseClient{
		ctx:    ctx,
		client: client,
	}, nil
}

// SetInvocationCount updates the invocation count and timestamp for a specific document identified by its ccna3 code.
func (client *FirebaseClient) SetInvocationCount(ccna3 string, number int) {
	// e.g. SetInvocationCount("NOR", 5)
	// finds the collection(CollectionInvocationCounts) -> create a reference to a document with document ID(ccna3)
	docRef := client.client.Collection(CollectionInvocationCounts).Doc(ccna3)
	_, err := docRef.Set(client.ctx, map[string]interface{}{
		"count": number,
		"time":  firestore.ServerTimestamp,
	})
	if err != nil {
		log.Printf("Failed to set invocation count: %v", err)
	}
}

// GetInvocationCount retrieves the invocation count for a given ccna3 (Country Code and Network Access Area)
func (client *FirebaseClient) GetInvocationCount(ccna3 string) (int, error) {
	// e.g. GetInvocationCount("NOR")

	// Get a reference to the document with the given ccna3 in the collection
	docRef := client.client.Collection(CollectionInvocationCounts).Doc(ccna3)

	// Fetch the document with the provided context
	docField, err := docRef.Get(client.ctx)
	if err != nil {
		// Log an error message if there was an issue retrieving the document
		log.Println("Error extracting body of returned document of message " + ccna3)
		return 0, err
	}
	// Get the value of the 'count' field from the document
	count, err := docField.DataAt("count")
	if err != nil {
		// Log an error message if there was an issue reading the 'count' field from the document
		log.Printf("Failed to read 'count' field from document: %v", err)
		return 0, err
	} else {
		// Log the invocation count for the given ccna3
		log.Printf("Invocation count for %s is %d", ccna3, count)
		switch c := count.(type) {
		case int:
			return c, nil
		case int64:
			return int(c), nil
		default:
			return 0, fmt.Errorf("unsupported type returned from GetInvocationCount: %T", c)
		}
	}
}

func (client *FirebaseClient) GetAllInvocationCounts() map[string]int {
	data := map[string]int{}
	docs, err := client.getAllDocuments(CollectionInvocationCounts)
	if err != nil {
		log.Printf("Could not fetch invocation counts from firestore")
		return data
	}
	for _, docField := range docs {
		if count, err := docField.DataAt("count"); err == nil {
			data[docField.Ref.ID] = int(count.(int64))
		}
	}
	return data
}

func (client *FirebaseClient) getAllDocuments(collection string) ([]*firestore.DocumentSnapshot, error) {
	docRef := client.client.Collection(collection).Documents(client.ctx)
	docs, err := docRef.GetAll()
	if err != nil {
		log.Printf("Failed to get documents: #{err}")
		return nil, err
	}
	return docs, nil
}

// SetRenewablesCache stores a YearRecordList in the renewables cache collection using the given URL as the document identifier.
func (client *FirebaseClient) SetRenewablesCache(url string, list *datastore.YearRecordList) {
	// e.g. SetRenewablesCache("/current/nor?neighbours=true", *data)
	// Access the renewables cache collection, with the specified url, if not exist, it will be created
	docRef := client.client.Collection(CollectionRenewablesCache).Doc(url)

	// Set or update the document with the provided data as a RenewableDB containing the YearRecordList
	_, err := docRef.Set(client.ctx, datastore.RenewableDB{
		"yearRecords": *list,
	})
	// Log errors if they occur during the Set operation
	if err != nil {
		log.Printf("Failed to set invocation registration: %v", err)
	}
}

// GetRenewablesCache retrieves a cached YearRecordList and its creation time by URL.
func (client *FirebaseClient) GetRenewablesCache(url string) (datastore.YearRecordList, time.Time, error) {
	// e.g. GetRenewablesCache("/current/nor?neighbours=true")
	// Access the renewables cache collection and get a reference to the document with the specified URL
	docRef := client.client.Collection(CollectionRenewablesCache).Doc(url)
	// Retrieve the document
	doc, err := docRef.Get(client.ctx)
	if err != nil {
		log.Printf("Failed to get renewables cache entry: %v", err)
		return datastore.YearRecordList{}, time.Time{}, err
	}
	// Extract the YearRecordList from the document
	var db datastore.RenewableDB
	err = doc.DataTo(&db)
	if err != nil {
		log.Printf("Failed to convert document data to RenewableDB: %v", err)
		return datastore.YearRecordList{}, time.Time{}, err
	}
	// Get the YearRecordList from the RenewableDB and creation time
	yearRecords := db["yearRecords"]
	creationTime := doc.CreateTime

	return yearRecords, creationTime, nil
}

// DeleteRenewablesCache removes the specified document from the renewables cache collection using the provided URL.
func (client *FirebaseClient) DeleteRenewablesCache(url string) {
	// e.g. DeleteRenewablesCache("/current/nor?neighbours=true")
	// Access the renewables cache collection and get a reference to the document with the specified URL
	docRef := client.client.Collection(CollectionRenewablesCache).Doc(url)
	// Delete the document
	_, err := docRef.Delete(client.ctx)
	// Log errors if they occur during the Delete operation
	if err != nil {
		log.Printf("Failed to delete renewables cache entry: %v", err)
	}
}

// SetInvocationRegistration stores an InvocationRegistration in Firestore
func (client *FirebaseClient) SetInvocationRegistration(registration InvocationRegistration) {
	// e.g. SetInvocationRegistration(registration)
	// Access the invocation_registrations collection, with the specified WebhookID, if not exist, it will be created
	docRef := client.client.Collection(CollectionInvocationRegistrations).Doc(registration.WebhookID)
	_, err := docRef.Set(client.ctx, registration) // Set or update the document with the provided data
	if err != nil {                                // Log errors if they occur during the Set operation
		log.Printf("Failed to set invocation registration: %v", err)
	}
}

// GetAllInvocationRegistrations retrieves all InvocationRegistration documents from Firestore
func (client *FirebaseClient) GetAllInvocationRegistrations() []InvocationRegistration {
	// e.g. GetAllInvocationRegistrations()
	// Create a DocumentIterator for the CollectionInvocationRegistrations collection
	docIterator := client.client.Collection(CollectionInvocationRegistrations).Documents(client.ctx)
	docs, err := docIterator.GetAll()
	if err != nil {
		log.Printf("Failed to get all invocation registrations: %v", err)
		return []InvocationRegistration{}
	}
	// Initialize an empty slice of InvocationRegistration
	registrations := []InvocationRegistration{}

	// Iterate through the retrieved documents and append them to the slice
	for _, doc := range docs {
		var registration InvocationRegistration
		if err := doc.DataTo(&registration); err != nil {
			log.Printf("Failed to convert document data to InvocationRegistration: %v", err)
			continue
		}
		registrations = append(registrations, registration)
	}
	return registrations
}

//func (client *FirebaseClient) BulkWriteInvocationCounts(ccna3map map[string]int) {
//	bulkWriter := client.client.BulkWriter(client.ctx)
//	for ccna3, count := range ccna3map {
//		docRef := client.client.Collection(CollectionInvocationCounts).Doc(ccna3)
//		bulkWriter.Set(docRef, map[string]interface{}{"count": count})
//	}
//	bulkWriter.End()
//
//}

func (client *FirebaseClient) BulkWrite(updates *BundledUpdate) {
	bulkWriter := client.client.BulkWriter(client.ctx)

	// updating invocation counts
	for ccna3, count := range updates.InvocationCount {
		docRef := client.client.Collection(CollectionInvocationCounts).Doc(ccna3)
		bulkWriter.Set(docRef, map[string]interface{}{"count": count})
	}

	// updating registrations
	for _, reg := range updates.Registrations {
		docRef := client.client.Collection(CollectionInvocationRegistrations).Doc(reg.Registration.WebhookID)
		if reg.Add {
			bulkWriter.Set(docRef, reg.Registration)
		} else {
			bulkWriter.Delete(docRef)
		}
	}

	// updating cache
	//for ccna3, count := range updates.Cache {
	//	docRef := client.client.Collection(CollectionInvocationCounts).Doc(ccna3)
	//	bulkWriter.Set(docRef, map[string]interface{}{"count": count})
	//}

	bulkWriter.End()

}
