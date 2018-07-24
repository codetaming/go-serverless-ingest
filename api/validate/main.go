package main

import (
	"encoding/json"
	"errors"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/codetaming/indy-ingest/api/utils"
	"github.com/codetaming/indy-ingest/api/validator"
	"log"
)

var (
	ErrMetadataNotProvided = errors.New("no metadata was provided in the HTTP body")
)

func Handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	log.Print(request)

	headers := map[string]string{"Content-Type": "application/json"}

	if len(request.Body) < 1 {
		return events.APIGatewayProxyResponse{}, ErrMetadataNotProvided
	}

	schemaUrl, err := utils.ExtractSchemaUrl(request.Headers)

	if err != nil {
		utils.RespondToClientError(err)
	}

	bodyJson := request.Body
	result, err := validator.Validate(schemaUrl, bodyJson)

	if err != nil {
		utils.RespondToClientError(err)
	}

	body, err := json.Marshal(result)

	if err != nil {
		utils.RespondToInternalError(err)
	}

	return events.APIGatewayProxyResponse{
		Headers:    headers,
		Body:       string(body),
		StatusCode: 200,
	}, nil

}

func main() {
	lambda.Start(Handler)
}
