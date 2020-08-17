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
	Values []struct {
		Value string `json:"value"`
	} `json:"values"`
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
	f, sheetURL, status, err := getResult(request, rType, format)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: status,
			Body:       `{ "error": "` + err.Error() + `" }`,
		}, nil
	}

	if f != nil {
		// If the return value is a file, serve it.
		return serveFile(f), nil
	} else {
		// If the return value is not a file, then it must be a sheetURL.
		return events.APIGatewayProxyResponse{
			StatusCode: status,
			Body:       `{ "sheetURL": "` + sheetURL + `" }`,
		}, nil
	}
}

// getResult returns the result by evaluating the option inputs.
func getResult(request events.APIGatewayProxyRequest, rType string, format string) (*bytes.Buffer, string, int, error) {
	// Unmarshal the json request.
	var rBody requestBody
	err := json.Unmarshal([]byte(request.Body), &rBody)
	if err != nil {
		return nil, "", http.StatusBadRequest, errors.New("Error occur while unmarshalling body-json value. Check your request.")
	}

	// Check the count.
	if len(rBody.Values) > 100 {
		return nil, "", http.StatusBadRequest, errors.New("You have more than 100 URLs.")
	}

	switch rType {
	case "url":
		switch format {
		case "excel":
			f, status, err := getExcelResultForURLs(rBody)
			return f, "", status, err
		case "sheet":
			sheetURL, status, err := getSheetResultForURLs(rBody)
			return nil, sheetURL, status, err
		default:
			return nil, "", http.StatusBadRequest, errors.New("Format must be \"excel\" or \"sheet\"")
		}

	case "keyword":
		switch format {
		case "excel":
			f, status, err := getExcelResultForKeywords(rBody)
			return f, "", status, err
		case "sheet":
			sheetURL, status, err := getSheetResultForKeywords(rBody)
			return nil, sheetURL, status, err
		default:
			return nil, "", http.StatusBadRequest, errors.New("Format must be \"excel\" or \"sheet\"")
		}
	default:
		return nil, "", http.StatusBadRequest, errors.New("Type must be \"url\" or \"keyword\".")
	}
}

// getExcelResultForURLs returns excel file for the given request.
func getExcelResultForURLs(rBody requestBody) (*bytes.Buffer, int, error) {
	// Create a new Set with inputs.
	urlSet := models.NewURLSet()
	for _, v := range rBody.Values {
		urlSet.Add(v.Value)
	}

	// Get the result
	status, err := services.GetResultFromSerpApiByUsingURLs(urlSet, "tr")
	if err != nil {
		return nil, status, err
	}

	// Convert the result to excel.
	f, err := services.ConvertURLResultToExcel(urlSet)
	if err != nil {
		return nil, http.StatusInternalServerError, errors.New("We have some issue while creating the excel output. Please try later.")
	}

	return f, http.StatusCreated, nil
}

// getSheetResultForURLs returns sheet url for the given request.
func getSheetResultForURLs(rBody requestBody) (string, int, error) {
	f, status, err := getExcelResultForURLs(rBody)
	if err != nil {
		return "", status, err
	}
	sheetURL, err := services.ImportFileToGoogleSheets(f)
	if err != nil {
		return "", http.StatusInternalServerError, errors.New("We have some issue with Google Sheets. Please try later.")
	}
	return sheetURL, http.StatusCreated, nil
}

// getExcelResultForKeywords returns excel file for the given request.
func getExcelResultForKeywords(rBody requestBody) (*bytes.Buffer, int, error) {
	// Create a new Set with inputs.
	keywordSet := models.NewKeywordSet()
	for _, v := range rBody.Values {
		keywordSet.Add(v.Value)
	}

	// Get the result
	status, err := services.GetResultFromSerpApiByUsingKeywords(keywordSet, "tr")
	if err != nil {
		return nil, status, err
	}

	// Convert the result to excel.
	f, err := services.ConvertKeywordResultToExcel(keywordSet)
	if err != nil {
		return nil, http.StatusInternalServerError, errors.New("We have some issue while creating the excel output. Please try later.")
	}

	return f, http.StatusCreated, nil
}

// getSheetResultForKeywords returns sheet url for the given request.
func getSheetResultForKeywords(rBody requestBody) (string, int, error) {
	f, status, err := getExcelResultForKeywords(rBody)
	if err != nil {
		return "", status, err
	}
	sheetURL, err := services.ImportFileToGoogleSheets(f)
	if err != nil {
		return "", http.StatusInternalServerError, errors.New("We have some issue with Google Sheets. Please try later.")
	}
	return sheetURL, http.StatusCreated, nil
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
