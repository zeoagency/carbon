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

func TestResult(t *testing.T) {
	// Check if method control is working.
	request := events.APIGatewayProxyRequest{}
	res, err := Result(request)
	if err == nil || res.StatusCode != http.StatusMethodNotAllowed {
		t.Fatal("Method checking is not working.")
	}

	request = events.APIGatewayProxyRequest{
		HTTPMethod: "POST",
		Body:       `{"urls": [{"url": "https://tools.zeo.org/carbon"}] }`,
	}

	res, err = Result(request)
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println(res)
}
