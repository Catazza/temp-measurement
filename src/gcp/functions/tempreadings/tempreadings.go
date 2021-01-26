package tempreadings

import (
	"fmt"
	// "context"
	// "log"
	"net/http"
	"time"
	// "cloud.google.com/go/bigquery"
)

// TODO: Put in shared module
type bqParsedTempTableRow struct {
	PubSubMessageID string    `bigquery:"pubsub_message_id"`
	ProcessingTime  time.Time `bigquery:"processing_time"`
	DeviceMessageID string    `json:"device_message_id" bigquery:"device_message_id"`
	Temperature     string    `json:"temperature" bigquery:"temperature"`
	Humidity        string    `json:"humidity" bigquery:"humidity"`
	MeasurementTime string    `json:"measurement_time" bigquery:"measurement_time"`
}

func RetrieveTempreadings(w http.ResponseWriter, r *http.Request) {
	// ctx := context.Background()

	// // TODO: REMOVE HARDCODING AND PUT ENV VAR
	// client, err := bigquery.NewClient(ctx, "temp-measure-dev")
	// if err != nil {
	// 	log.Fatalf("Error initialising client %v", err)
	// }
	// defer client.Close()

	// // TODO: REMOVE HARDCODING AND PUT ENV VAR
	// row, err := client.Query(
	// 	`SELECT * FROM ` + "`temp-measure-dev.temp_measure.temp_history_parsed`" + ` WHERE DATE(processing_time) = DATE(CURRENT_DATETIME())
	// 	order by measurement_time desc
	// 	LIMIT 100`).Read(ctx)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	fmt.Fprint(w, "hola chicaa")

}
