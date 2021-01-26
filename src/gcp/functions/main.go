package main

import (
	"context"
	"log"
	"os"

	"azzadigital.com/tempmeasurement/cloudfunctions/dbloader"
	"azzadigital.com/tempmeasurement/cloudfunctions/tempreadings"
	"github.com/GoogleCloudPlatform/functions-framework-go/funcframework"
)

func main() {
	ctx := context.Background()
	// m := metadata.Metadata{EventID: uuid.New().String(), Timestamp: time.Now()}
	// ctxWithMetadata := metadata.NewContext(ctx, &m)

	// t, err := metadata.FromContext(ctxWithMetadata)
	// log.Println(t, err)

	if err := funcframework.RegisterEventFunctionContext(ctx, "/dbloader", dbloader.StoreTempMeasurementBQ); err != nil {
		log.Fatalf("funcframework.RegisterEventFunctionContext: %v\n", err)
	}

	if err := funcframework.RegisterHTTPFunctionContext(ctx, "/tempreadings", tempreadings.RetrieveTempreadings); err != nil {
		log.Fatalf("funcframework.RegisterEventFunctionContext: %v\n", err)
	}

	// Use PORT environment variable, or default to 8080.
	port := "8080"
	if envPort := os.Getenv("PORT"); envPort != "" {
		port = envPort
	}
	if err := funcframework.Start(port); err != nil {
		log.Fatalf("funcframework.Start: %v\n", err)
	}
}
