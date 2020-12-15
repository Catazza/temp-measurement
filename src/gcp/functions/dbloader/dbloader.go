package dbloader

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"cloud.google.com/go/bigquery"
	"google.golang.org/api/iterator"
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
	DBLoaderProcessingTime time.Time
}

type StackOverflowRow struct {
	URL       string `bigquery:"url"`
	ViewCount int64  `bigquery:"view_count"`
}

// printResults prints results from a query to the Stack Overflow public dataset.
func printResults(w io.Writer, iter *bigquery.RowIterator) error {
	for {
		var row StackOverflowRow
		err := iter.Next(&row)
		if err == iterator.Done {
			return nil
		}
		if err != nil {
			return fmt.Errorf("error iterating through results: %v", err)
		}

		fmt.Fprintf(w, "url: %s views: %d\n", row.URL, row.ViewCount)
	}
}

// StoreTempMeasurementBQ consumes a Pub/Sub message.
func StoreTempMeasurementBQ(ctx context.Context, m PubSubMessage) error {
	pubsubPayload := string(m.Data) // Automatically decoded from base64.
	tempMeasurement := TempMeasurement{}
	json.Unmarshal([]byte(pubsubPayload), &tempMeasurement)
	log.Printf("Payload: %s", pubsubPayload)

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

	// example to check if reading works
	query := client.Query(
		`SELECT
                CONCAT(
                        'https://stackoverflow.com/questions/',
                        CAST(id as STRING)) as url,
                view_count
        FROM ` + "`bigquery-public-data.stackoverflow.posts_questions`" + `
        WHERE tags like '%google-bigquery%'
        ORDER BY view_count DESC
        LIMIT 10;`)

	results, err := query.Read(ctxWithTimeout)
	if err != nil {
		log.Fatal(err)
	}

	if err := printResults(os.Stdout, results); err != nil {
		log.Fatal(err)
	}

	return nil
}
