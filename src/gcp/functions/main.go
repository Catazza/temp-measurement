package main

import (
	"context"
	"log"
	"os"

	"azzadigital.com/tempmeasurement/cloudfunctions/dbloader"
	"github.com/GoogleCloudPlatform/functions-framework-go/funcframework"
)

func main() {
	ctx := context.Background()
	// if err := funcframework.RegisterHTTPFunctionContext(ctx, "/", hello.HelloWorld); err != nil {
	// log.Fatalf("funcframework.RegisterHTTPFunctionContext: %v\n", err)
	// }

	if err := funcframework.RegisterEventFunctionContext(ctx, "/", dbloader.StoreTempMeasurementBQ); err != nil {
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
