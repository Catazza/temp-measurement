package dbloader

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"cloud.google.com/go/bigquery"
	"cloud.google.com/go/functions/metadata"
)

// PubSubMessage is the payload of a Pub/Sub event.
type PubSubMessage struct {
	Data []byte `json:"data"`
}

type bqRawTempTableRow struct {
	PubSubMessageID string    `bigquery:"pubsub_message_id"`
	JSONMsg         string    `bigquery:"json_msg"`
	ProcessingTime  time.Time `bigquery:"processing_time"`
}

type bqParsedTempTableRow struct {
	PubSubMessageID string    `bigquery:"pubsub_message_id"`
	ProcessingTime  time.Time `bigquery:"processing_time"`
	DeviceMessageID string    `json:"device_message_id" bigquery:"device_message_id"`
	Temperature     string    `json:"temperature" bigquery:"temperature"`
	Humidity        string    `json:"humidity" bigquery:"humidity"`
	MeasurementTime string    `json:"measurement_time" bigquery:"measurement_time"`
}

func retrievePubSubMessageID(ctx context.Context) string {
	// Retrive pubSub message ID
	ctxMetadata, err := metadata.FromContext(ctx)
	if err != nil && err.Error() == "unable to find metadata" {
		// shortcut to test locally, to be improved to add local metadata in context
		return "99999" // "random" ID
	}
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("MSG ID: %s", ctxMetadata.EventID)
	return ctxMetadata.EventID
}

// StoreTempMeasurementBQ consumes a Pub/Sub message.
func StoreTempMeasurementBQ(ctx context.Context, m PubSubMessage) error {
	// retrieve the measurement payload
	pubsubPayload := string(m.Data) // Automatically decoded from base64.

	// retrieve and compute attributes to set, such as the PS message ID
	msgID := retrievePubSubMessageID(ctx)
	processingTime := time.Now()

	// Create BQ items to save
	bqRawItem := bqRawTempTableRow{PubSubMessageID: msgID, JSONMsg: pubsubPayload, ProcessingTime: processingTime}
	bqParsedItem := bqParsedTempTableRow{}
	json.Unmarshal([]byte(pubsubPayload), &bqParsedItem)
	bqParsedItem.PubSubMessageID = msgID
	bqParsedItem.ProcessingTime = processingTime

	// Initialise BQ client
	ctxWithTimeout, cancelFunction := context.WithTimeout(ctx, time.Duration(200)*time.Second)
	client, err := bigquery.NewClient(ctxWithTimeout, "temp-measure-dev")
	if err != nil {
		log.Fatalf("bigquery.Newclient: %v", err)
	}

	defer func() {
		fmt.Println("BQ Loader returned - canceling context")
		client.Close()
		cancelFunction()
	}()

	// Save objects to Bigquery
	// Raw line
	inserter := client.Dataset("temp_measure").Table("temp_history_raw").Inserter()
	saveErr := inserter.Put(ctx, &bqRawItem)
	if saveErr != nil {
		log.Fatalf("Error saving item %v to BQ table %v: %v", &bqRawItem, "temp_history_raw", saveErr)
	}

	// Parsed Line
	inserterP := client.Dataset("temp_measure").Table("temp_history_parsed").Inserter()
	saveErrP := inserterP.Put(ctx, &bqParsedItem)
	if saveErrP != nil {
		log.Fatalf("Error saving item %v to BQ table %v: %v", &bqParsedItem, "temp_history_parsed", saveErrP)
	}
	return nil
}
