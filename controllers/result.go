package controllers

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/aws/aws-lambda-go/events"
	"gitlab.com/seo.do/zeo-carbon/models"
	"gitlab.com/seo.do/zeo-carbon/services"
)

// requestBody keeps request body.
type requestBody struct {
	Urls []struct {
		URL string `json:"url"`
	} `json:"urls"`
}

// Result works like router.
//
// You need to send type and format in the request.
// You will get a response that is related with request.
func Result(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// Get params, returns an error if the param is not set.
	rType, format, status, err := checkAndGetParams(request)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: status,
			Body:       `{ "error": "` + err.Error() + `" }`,
		}, nil
	}

	// Process the request.
	switch rType {
	case "url":
		switch format {
		case "excel":
			f, status, err := getExcelResultForURLs(request)
			if err != nil {
				return events.APIGatewayProxyResponse{
					StatusCode: status,
					Body:       `{ "error": "` + err.Error() + `" }`,
				}, nil
			}

			// Serve the result.
			return serveFile(f), nil
		case "sheet":
			return events.APIGatewayProxyResponse{
				StatusCode: http.StatusServiceUnavailable,
				Body:       `{ "error": "` + "Sheet format is not available yet. Please try later." + `" }`,
			}, nil
		default:
			return events.APIGatewayProxyResponse{
				StatusCode: http.StatusBadRequest,
				Body:       `{ "error": "` + "Format must be \"excel\" or \"sheet\"." + `" }`,
			}, nil
		}
	case "keyword":
		switch format {
		case "excel":
			f, status, err := getExcelResultForKeywords(request)
			if err != nil {
				return events.APIGatewayProxyResponse{
					StatusCode: status,
					Body:       `{ "error": "` + err.Error() + `" }`,
				}, nil
			}

			// Serve the result.
			return serveFile(f), nil
		case "sheet":
			return events.APIGatewayProxyResponse{
				StatusCode: http.StatusServiceUnavailable,
				Body:       `{ "error": "` + "Sheet format is not available yet. Please try later." + `" }`,
			}, nil
		default:
			return events.APIGatewayProxyResponse{
				StatusCode: http.StatusBadRequest,
				Body:       `{ "error": "` + "Format must be \"excel\" or \"sheet\"." + `" }`,
			}, nil
		}
	default:
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusBadRequest,
			Body:       `{ "error": "` + "Type must be \"url\" or \"keyword\"." + `" }`,
		}, nil
	}

}

// getExcelResultForURLs returns excel file for the given request.
func getExcelResultForURLs(request events.APIGatewayProxyRequest) (*bytes.Buffer, int, error) {
	// Unmarshal the json request.
	var br requestBody
	err := json.Unmarshal([]byte(request.Body), &br)
	if err != nil {
		return nil, http.StatusBadRequest, errors.New("Error occur while unmarshalling body-json value. Check your request.")
	}

	// Check the count.
	if len(br.Urls) > 100 {
		return nil, http.StatusBadRequest, errors.New("You have more than 100 URLs.")
	}

	// Create a new Set with inputs.
	urlSet := models.NewURLSet()
	for _, v := range br.Urls {
		err := urlSet.Add(v.URL)
		if err != nil {
			return nil, http.StatusBadRequest, err
		}
	}

	// Get the result
	success, fail, status, err := services.GetResultFromSerpApiByUsingURLs(urlSet, "tr")
	if err != nil {
		return nil, status, err
	}

	// Convert the result to excel.
	f, err := services.ConvertURLResultToExcel(success, fail)
	if err != nil {
		return nil, http.StatusInternalServerError, errors.New("We have some issue while creating the excel output. Please try later.")
	}

	return f, http.StatusCreated, nil
}

// getExcelResultForKeywords returns excel file for the given request.
func getExcelResultForKeywords(request events.APIGatewayProxyRequest) (*bytes.Buffer, int, error) {
	return nil, http.StatusServiceUnavailable, errors.New("Keyword endpoint is not available yet. Please try later.")
}

// checkAndGetParams checks the params are set or not.
func checkAndGetParams(request events.APIGatewayProxyRequest) (string, string, int, error) {
	// Check the method.
	if request.HTTPMethod != "POST" {
		return "", "", http.StatusMethodNotAllowed, errors.New("Method not allowed. Only allowed: POST.")
	}

	// Check the type.
	rType := ""
	if v, ok := request.QueryStringParameters["type"]; ok {
		rType = v
	} else {
		return "", "", http.StatusBadRequest, errors.New("Type is not set.")
	}

	// Check the format.
	format := ""
	if v, ok := request.QueryStringParameters["format"]; ok {
		format = v
	} else {
		return "", "", http.StatusBadRequest, errors.New("Format is not set.")
	}

	return rType, format, http.StatusOK, nil
}

// serveFile create a response to serve the given file.
func serveFile(f *bytes.Buffer) events.APIGatewayProxyResponse {
	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusCreated,
		Headers: map[string]string{
			"Content-Disposition": `attachment; filename="result.xlsx"`,
			"Content-Length":      strconv.Itoa(f.Len()),
			"Content-Type":        "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet",
		},
		Body:            base64.StdEncoding.EncodeToString(f.Bytes()),
		IsBase64Encoded: true,
	}
}
