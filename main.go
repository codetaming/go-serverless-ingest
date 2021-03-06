package main

import (
	"github.com/codetaming/indy-ingest/cmd/api"
	"github.com/codetaming/indy-ingest/internal/persistence/aws"
	"github.com/codetaming/indy-ingest/internal/validator"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
)

var (
	serverPort      = os.Getenv("SERVER_PORT")
	region          = os.Getenv("AWS_REGION")
	datasetTable    = os.Getenv("DATASET_TABLE")
	metadataTable   = os.Getenv("METADATA_TABLE")
	metadataBucket  = os.Getenv("METADATA_BUCKET")
	schemaCacheFile = os.Getenv("SCHEMA_CACHE_FILE")
)

func init() {
	if serverPort == "" {
		log.Fatal("$SERVER_PORT not set")
	}
	if region == "" {
		log.Fatal("$AWS_REGION not set")
	}
	if datasetTable == "" {
		log.Fatal("$DATASET_TABLE not set")
	}
	if metadataTable == "" {
		log.Fatal("$METADATA_TABLE not set")
	}
	if metadataBucket == "" {
		log.Fatal("$METADATA_BUCKET not set")
	}
	if metadataBucket == "" {
		log.Fatal("$METADATA_BUCKET not set")
	}
	if schemaCacheFile == "" {
		log.Fatal("$SCHEMA_CACHE_FILE not set")
	}
}

func main() {
	router := mux.NewRouter()

	logger := log.New(os.Stdout, "ingest ", log.LstdFlags|log.Lshortfile)

	logger.Printf("starting ingest")

	logger.Printf("configuring data store")
	dataStore, err := aws.NewDynamoDataStore(logger, region, datasetTable, metadataTable)
	if err != nil {
		logger.Fatalf("failed to create data store: %v", err)
	}

	logger.Printf("configuring file store")
	fileStore, err := aws.NewS3FileStore(logger, region, metadataBucket)
	if err != nil {
		logger.Fatalf("failed to create file store: %v", err)
	}

	logger.Printf("configuring validator")
	validator, err := validator.NewCachingValidator(logger, schemaCacheFile)
	if err != nil {
		logger.Fatalf("failed to create validator: %v", err)
	}

	a := api.NewAPI(logger, dataStore, fileStore, validator)
	a.SetupRoutes(router)

	logger.Printf("server starting on port %s", serverPort)
	err = http.ListenAndServe(":"+serverPort, router)
	if err != nil {
		logger.Fatalf("server failed to start: %v", err)
	}
}
