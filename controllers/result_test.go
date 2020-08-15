package controllers

import (
	"fmt"
	"log"
	"net/http"
	"testing"

	"github.com/aws/aws-lambda-go/events"
	"github.com/joho/godotenv"
)

func init() {
	err := godotenv.Load("../.env")
	if err != nil {
		log.Fatalln(err)
	}
}

func TestResultShouldFail(t *testing.T) {
	// Check if method control is working.
	request := events.APIGatewayProxyRequest{}
	res, _ := Result(request)
	if res.StatusCode != http.StatusMethodNotAllowed {
		t.Fatal("Method checking is not working.")
	}
}

func TestResultForURLs(t *testing.T) {
	request := events.APIGatewayProxyRequest{
		HTTPMethod: "POST",
		QueryStringParameters: map[string]string{
			"type":   "url",
			"format": "excel",
		},
		Body: `{"urls": [{"url": "https://tools.zeo.org/carbon"}] }`,
	}

	res, _ := Result(request)
	if res.StatusCode != http.StatusCreated {
		t.Fatal("Error occur while getting excel file.")
	}

	fmt.Println(res.Headers)
}
