package dbloader

import (
	"context"
	"encoding/json"
	"log"
)

// struct to puss in the cloud function
type PubSubMessage struct {
	Data []byte
	// publishTime tim
}

type TemperaturePayload struct {
	temperature   string
	humidity      string
	measuringTime string
}

func DBLoader(ctx context.Context, m PubSubMessage) error {
	log.Printf("data before marshalling: %s", string(m.Data))
	var parsedTempPayload TemperaturePayload
	json.Unmarshal(m.Data, &parsedTempPayload)
	log.Printf("data is %s", parsedTempPayload)
	return nil
}
