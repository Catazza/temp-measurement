package dbloader

import (
	"context"
	"encoding/json"
	"log"
	"time"
)

// PubSubMessage is the payload of a Pub/Sub event.
type PubSubMessage struct {
	Data []byte `json:"data"`
}

// TempMeasurement is the struct to Unmarshal the temperature measurement
type TempMeasurement struct {
	Temperature   string    `json:"temperature"`
	Humidity      string    `json:"humidity"`
	MeasuringTime time.Time `json:"measuring_time"`
}

// StoreTempMeasurementBQ consumes a Pub/Sub message.
func StoreTempMeasurementBQ(ctx context.Context, m PubSubMessage) error {
	pubsubPayload := string(m.Data) // Automatically decoded from base64.
	tempMeasurement := TempMeasurement{}
	json.Unmarshal([]byte(pubsubPayload), &tempMeasurement)
	log.Printf("Payload: %s!", pubsubPayload)
	return nil
}
