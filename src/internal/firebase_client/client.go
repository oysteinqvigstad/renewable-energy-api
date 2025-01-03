package firebase_client

import (
	"assignment2/internal/types"
	"cloud.google.com/go/firestore" // Firestore-specific support
	"context"                       // State handling across API boundaries; part of native GoLang API
	firebase "firebase.google.com/go"
	"google.golang.org/api/option"
	"log"
	"strings"
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
	ctx := context.Background()                                  // Create a basic empty box for tasks
	secretKeyPath := option.WithCredentialsFile(PathToSecretKey) // Tell the program where to find the secret key for Firebase
	app, err := firebase.NewApp(ctx, nil, secretKeyPath)         // Connect to Firebase using the secret key

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
func (client *FirebaseClient) GetInvocationCount(ccna3 string) (int64, error) {
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
	}
	return count.(int64), nil
}

func (client *FirebaseClient) GetAllInvocationCounts() map[string]int64 {
	data := map[string]int64{}
	docs, err := client.GetAllDocuments(CollectionInvocationCounts)
	if err != nil {
		log.Printf("Could not fetch invocation counts from firestore")
		return data
	}
	for _, docField := range docs {
		if count, err := docField.DataAt("count"); err == nil {
			data[docField.Ref.ID] = count.(int64)
		}
	}
	return data
}

func (client *FirebaseClient) GetAllDocuments(collection string) ([]*firestore.DocumentSnapshot, error) {
	docRef := client.client.Collection(collection).Documents(client.ctx)
	docs, err := docRef.GetAll()
	if err != nil {
		log.Printf("Failed to get documents: #{err}")
		return nil, err
	}
	return docs, nil
}

// SetRenewablesCache stores a YearRecordList in the renewables cache collection using the given URL as the document identifier.
func (client *FirebaseClient) SetRenewablesCache(url string, list types.YearRecordList) {
	url = strings.ReplaceAll(url, "/", "_")
	// e.g. SetRenewablesCache("/current/nor?neighbours=true", *data)
	// Access the renewables cache collection, with the specified url, if not exist, it will be created
	docRef := client.client.Collection(CollectionRenewablesCache).Doc(url)

	// Set or update the document with the provided data as a RenewableDB containing the YearRecordList
	_, err := docRef.Set(client.ctx, types.RenewableDB{
		"yearRecords": list,
	})
	// Log errors if they occur during the Set operation
	if err != nil {
		log.Printf("Failed to set invocation registration: %v", err)
	}
}

// GetRenewablesCache retrieves a cached YearRecordList and its creation time by URL.
func (client *FirebaseClient) GetRenewablesCache(url string) (types.YearRecordList, time.Time, error) {
	url = strings.ReplaceAll(url, "/", "_")
	// Access the renewables cache collection and get a reference to the document with the specified URL
	docRef := client.client.Collection(CollectionRenewablesCache).Doc(url)
	// Retrieve the document
	doc, err := docRef.Get(client.ctx)
	if err != nil {
		log.Printf("Failed to get renewables cache entry: %v", err)
		return types.YearRecordList{}, time.Time{}, err
	}
	// Extract the YearRecordList from the document
	var db types.RenewableDB
	err = doc.DataTo(&db)
	if err != nil {
		log.Printf("Failed to convert document data to RenewableDB: %v", err)
		return types.YearRecordList{}, time.Time{}, err
	}
	// Get the YearRecordList from the RenewableDB and creation time
	yearRecords := db["yearRecords"]
	creationTime := doc.CreateTime

	return yearRecords, creationTime, nil
}

// DeleteRenewablesCache removes the specified document from the renewables cache collection using the provided URL.
func (client *FirebaseClient) DeleteRenewablesCache(url string) {
	url = strings.ReplaceAll(url, "/", "_")
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
func (client *FirebaseClient) SetInvocationRegistration(registration types.InvocationRegistration) {
	// Access the invocation_registrations collection, with the specified WebhookID, if not exist, it will be created
	docRef := client.client.Collection(CollectionInvocationRegistrations).Doc(registration.WebhookID)
	_, err := docRef.Set(client.ctx, registration) // Set or update the document with the provided data
	if err != nil {                                // Log errors if they occur during the Set operation
		log.Printf("Failed to set invocation registration: %v", err)
	}
}

// GetAllInvocationRegistrations retrieves all InvocationRegistration documents from Firestore
func (client *FirebaseClient) GetAllInvocationRegistrations() map[string]types.InvocationRegistration {
	result := map[string]types.InvocationRegistration{}
	docs, err := client.GetAllDocuments(CollectionInvocationRegistrations)
	if err != nil {
		log.Printf("Could not fetch data from firestore")
		return result
	}
	for _, doc := range docs {
		var registration types.InvocationRegistration
		err = doc.DataTo(&registration)
		if err != nil {
			return result
		}
		result[doc.Ref.ID] = registration

	}
	return result
}

func (client *FirebaseClient) BulkWrite(updates *types.BundledUpdate) {
	bulkWriter := client.client.BulkWriter(client.ctx)

	// updating invocation counts
	for countryCode, count := range updates.InvocationCount {
		docRef := client.client.Collection(CollectionInvocationCounts).Doc(countryCode)
		_, err := bulkWriter.Set(docRef, map[string]interface{}{"count": count})
		if err != nil {
			log.Println("could not add job to bulk-writer ", err.Error())
		}
	}

	// updating registrations
	for _, reg := range updates.Registrations {
		docRef := client.client.Collection(CollectionInvocationRegistrations).Doc(reg.Registration.WebhookID)
		if reg.Add {
			_, err := bulkWriter.Set(docRef, reg.Registration)
			if err != nil {
				log.Println("could not add job to bulk-writer ", err.Error())
			}
		} else {
			_, err := bulkWriter.Delete(docRef)
			if err != nil {
				log.Println("could not add job to bulk-writer ", err.Error())
			}
		}
	}

	// updating cache
	for url, cache := range updates.Cache {
		url = strings.ReplaceAll(url, "/", "_")
		docRef := client.client.Collection(CollectionRenewablesCache).Doc(url)
		_, err := bulkWriter.Set(docRef, map[string]interface{}{"yearRecords": cache})
		if err != nil {
			log.Println("could not add job to bulk-writer ", err.Error())
		}
	}

	bulkWriter.End()
}

func (client *FirebaseClient) Close() {
	_ = client.client.Close()

}
