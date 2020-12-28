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

// TempMeasurement is the struct to Unmarshal the temperature measurement
type TempMeasurement struct {
	Temperature            string    `json:"temperature"`
	Humidity               string    `json:"humidity"`
	MeasuringTime          time.Time `json:"measuring_time"`
	DeviceMessageID        string    `json:"device_message_id"`
	PubSubMessageID        string
	DBLoaderProcessingTime time.Time
}

type bqRawTempTableRow struct {
	PubSubMessageID string    `bigquery:"pubsub_message_id"`
	JSONMsg         string    `bigquery:"json_msg"`
	ProcessingTime  time.Time `bigquery:"processing_time"`
}

func retrievePubSubMessageID(ctx context.Context) string {
	// Retrive pubSub message ID
	ctxMetadata, err := metadata.FromContext(ctx)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("MSG ID: %s", ctxMetadata.EventID)
	return ctxMetadata.EventID
}

// StoreTempMeasurementBQ consumes a Pub/Sub message.
func StoreTempMeasurementBQ(ctx context.Context, m PubSubMessage) error {
	pubsubPayload := string(m.Data) // Automatically decoded from base64.
	tempMeasurement := TempMeasurement{}
	json.Unmarshal([]byte(pubsubPayload), &tempMeasurement)
	log.Printf("Payload: %s", pubsubPayload)

	// retrieve PS message ID
	// msgID := retrievePubSubMessageID(ctx)

	// Create BQ item to save
	bqRawItem := bqRawTempTableRow{PubSubMessageID: "11213", JSONMsg: pubsubPayload, ProcessingTime: time.Now()}

	// Initialise BQ client
	ctxWithTimeout, cancelFunction := context.WithTimeout(ctx, time.Duration(15)*time.Second)
	client, err := bigquery.NewClient(ctxWithTimeout, "temp-measure-dev")
	if err != nil {
		log.Fatalf("bigquery.Newclient: %v", err)
	}

	defer func() {
		fmt.Println("BQ Loader returned - canceling context")
		client.Close()
		cancelFunction()
	}()

	// Save object to Bigquery
	inserter := client.Dataset("temp_measure").Table("temp_history_raw").Inserter()
	saveErr := inserter.Put(ctxWithTimeout, &bqRawItem)
	if saveErr != nil {
		log.Fatalf("Error saving item %v to BQ: %v", &bqRawItem, saveErr)
	}

	return nil
}
