package main

import (
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/codetaming/indy-ingest/api/model"
	"github.com/codetaming/indy-ingest/api/persistence"
	"github.com/google/uuid"
	"os"
	"time"
)

//AWS Lambda entry point
func Handler(_ events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return Do(new(persistence.DynamoPersistence))
}

//Do executes the function allowing dependencies to be specified
func Do(p persistence.DatasetPersister) (events.APIGatewayProxyResponse, error) {
	return respond(createDataSet(p))
}

func createDataSet(p persistence.DatasetPersister) (model.Dataset, error) {
	d := model.Dataset{
		Owner:     model.DefaultOwner,
		DatasetId: uuid.Must(uuid.NewUUID()).String(),
		Created:   time.Now(),
	}
	return d, p.PersistDataset(d)
}

func respond(d model.Dataset, err error) (events.APIGatewayProxyResponse, error) {
	headers := map[string]string{"Content-Type": "application/json"}
	if err != nil {
		return events.APIGatewayProxyResponse{
			Headers:    headers,
			Body:       err.Error(),
			StatusCode: 500,
		}, nil
	}
	baseUrl := os.Getenv("BASE_URL")
	headers["Location"] = baseUrl + "/dataset/" + d.DatasetId
	body, _ := json.Marshal(d)
	return events.APIGatewayProxyResponse{
		Headers:    headers,
		Body:       string(body),
		StatusCode: 201,
	}, nil
}

func main() {
	lambda.Start(Handler)
}
