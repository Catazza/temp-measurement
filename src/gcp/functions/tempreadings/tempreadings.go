package tempreadings

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"cloud.google.com/go/bigquery"
	"google.golang.org/api/iterator"
)

// TODO: Put in shared module
type bqParsedTempTableRow struct {
	PubSubMessageID string    `bigquery:"pubsub_message_id"`
	ProcessingTime  time.Time `bigquery:"processing_time"`
	DeviceMessageID string    `json:"device_message_id" bigquery:"device_message_id"`
	Temperature     float32   `json:"temperature" bigquery:"temperature"`
	Humidity        float32   `json:"humidity" bigquery:"humidity"`
	MeasurementTime time.Time `json:"measurement_time" bigquery:"measurement_time"`
}

func RetrieveTempreadings(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	// TODO: REMOVE HARDCODING AND PUT ENV VAR
	client, err := bigquery.NewClient(ctx, "temp-measure-dev")
	if err != nil {
		log.Fatalf("Error initialising client %v", err)
	}
	defer client.Close()

	// TODO: REMOVE HARDCODING AND PUT ENV VAR
	rows, err := client.Query(
		`SELECT * FROM ` + "`temp-measure-dev.temp_measure.temp_history_parsed`" + ` WHERE DATE(processing_time) = DATE(CURRENT_DATETIME())
		order by measurement_time desc
		LIMIT 100`).Read(ctx)
	if err != nil {
		log.Fatal(err)
	}

	// Initialise container
	tempMeasurements := make([]bqParsedTempTableRow, 0)

	for {
		var singleRow bqParsedTempTableRow
		err := rows.Next(&singleRow)
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Fatalf("Error in parsing row: %v; %v", singleRow, err)
		}
		tempMeasurements = append(tempMeasurements, singleRow)
	}
	w.Header().Set("Content-Type", "application/json")
	// fmt.Fprintf(w, tempMeasurements)
	json.NewEncoder(w).Encode(tempMeasurements)
}
