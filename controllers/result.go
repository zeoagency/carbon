package controllers

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/aws/aws-lambda-go/events"

	"gitlab.com/seo.do/zeo-carbon/models"
	"gitlab.com/seo.do/zeo-carbon/services"
)

type bodyRequest struct {
	Urls []struct {
		URL string `json:"url"`
	} `json:"urls"`
}

func Result(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// Check method.
	if request.HTTPMethod != "POST" {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusMethodNotAllowed,
		}, errors.New("Method not allowed. Only allowed: POST.")
	}

	// Unmarshal the json request.
	var br bodyRequest
	err := json.Unmarshal([]byte(request.Body), &br)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusBadRequest,
		}, errors.New("Error occur while unmarshalling body-json value. Check your request.")
	}

	// Check the count.
	if len(br.Urls) > 100 {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusBadRequest,
		}, errors.New("You have more than 100 URLs.")
	}

	// Create a new Set with inputs.
	urlSet := models.NewURLSet()
	for _, v := range br.Urls {
		err := urlSet.Add(v.URL)
		if err != nil {
			return events.APIGatewayProxyResponse{
				StatusCode: http.StatusBadRequest,
			}, errors.New("You entred non-url input. Check your URLS.")
		}
	}

	// Get the result
	result, status, err := services.GetResultFromSerpApiByUsingURLs(urlSet, "tr")
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: status,
		}, errors.New(err.Error())
	}

	// Marshall the result.
	r, err := json.MarshalIndent(result, "", "    ")
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
		}, errors.New("JSON marshall error occur.")
	}

	// Serve the result.
	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusCreated,
		Body:       string(r),
	}, nil
}
