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

func TestExcelResultForURLs(t *testing.T) {
	request := events.APIGatewayProxyRequest{
		HTTPMethod: "POST",
		QueryStringParameters: map[string]string{
			"type":     "url",
			"format":   "excel",
			"country":  "tr",
			"language": "tr",
		},
		Body: `{"values": [{"value": "https://tools.zeo.org/carbon"}, {"value": "https://zeo.org"}] }`,
	}

	res, _ := Result(request)
	if res.StatusCode != http.StatusCreated {
		t.Fatal("Error occur while getting excel file.", res.Body)
	}

	fmt.Println(res.Headers)
}

func TestExcelResultForKeywords(t *testing.T) {
	request := events.APIGatewayProxyRequest{
		HTTPMethod: "POST",
		QueryStringParameters: map[string]string{
			"type":            "keyword",
			"format":          "excel",
			"country":         "tr",
			"language":        "tr",
			"accountName":     "bora@zeo.org",
			"accountPassword": "12341234",
		},
		Body: `{"values": [{"value": "zeo carbon tool"}] }`,
	}

	res, _ := Result(request)
	if res.StatusCode != http.StatusCreated {
		t.Fatal("Error occur while getting excel file.", res.Body)
	}

	fmt.Println(res.Headers)
}

func TestSheetResultForURLs(t *testing.T) {
	request := events.APIGatewayProxyRequest{
		HTTPMethod: "POST",
		QueryStringParameters: map[string]string{
			"type":     "url",
			"format":   "sheet",
			"country":  "tr",
			"language": "tr",
		},
		Body: `{"values": [{"value": "https://tools.zeo.org/carbon"}] }`,
	}

	res, _ := Result(request)
	if res.StatusCode != http.StatusCreated {
		t.Fatal("Error occur while getting excel file.", res.Body)
	}

	fmt.Println(res.Body)
}

func TestSheetResultForKeywords(t *testing.T) {
	request := events.APIGatewayProxyRequest{
		HTTPMethod: "POST",
		QueryStringParameters: map[string]string{
			"type":            "keyword",
			"format":          "sheet",
			"country":         "tr",
			"language":        "tr",
			"accountName":     "bora@zeo.org",
			"accountPassword": "12341234",
		},
		Body: `{"values": [{"value": "zeo carbon tool"}] }`,
	}

	res, _ := Result(request)
	if res.StatusCode != http.StatusCreated {
		t.Fatal("Error occur while getting excel file.", res.Body)
	}

	fmt.Println(res.Body)
}
