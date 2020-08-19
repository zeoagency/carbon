package controllers

import (
	"bytes"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/aws/aws-lambda-go/events"

	"github.com/zeoagency/carbon/models"
	"github.com/zeoagency/carbon/services"
)

// requestBody keeps request body.
type requestBody struct {
	Values []struct {
		Value string `json:"value"`
	} `json:"values"`
}

type internal struct {
	Accounts []struct {
		Name     string `json:"name"`
		Password string `json:"password"`
		Limit    int    `json:"limit"` // "-1" means there is no limit.
	} `json:"accounts"`
}

var rType, format, country, language string

// Result works like router.
//
// You need to send type and format in the request.
// You will get a response that is related with request.
func Result(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// Set params, returns an error if the param is not set.
	status, err := checkAndSetParams(request)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: status,
			Body:       `{ "error": "` + err.Error() + `" }`,
		}, nil
	}

	// Check internal.
	// iLimit
	//    "0" = non-login user.
	//    "-1" = limitless user.
	//    "..." = limit is defined.
	isInternal, iLimit, status, err := checkAndAuthInternal(request)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: status,
			Body:       `{ "error": "` + err.Error() + `" }`,
		}, nil
	}

	// Process the request.
	f, sheetURL, status, err := getResult(request, isInternal, iLimit)
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
func getResult(request events.APIGatewayProxyRequest, isInternal bool, iLimit int) (*bytes.Buffer, string, int, error) {
	// Unmarshal the json request.
	var rBody requestBody
	err := json.Unmarshal([]byte(request.Body), &rBody)
	if err != nil {
		return nil, "", http.StatusBadRequest, errors.New("Error occur while unmarshalling body-json value. Check your request.")
	}

	status, err := checkLimit(len(rBody.Values), isInternal, iLimit)
	if err != nil {
		return nil, "", status, err
	}

	if rType == "url" {
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
	} else if isInternal && rType == "keyword" {
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
	} else {
		errText := "Type must be \"url\"."
		if isInternal {
			errText = "Type must be \"url\" or \"keyword\"."
		}
		return nil, "", http.StatusBadRequest, errors.New(errText)
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
	status, err := services.GetResultFromSerpApiByUsingURLs(urlSet, country, language)
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
	status, err := services.GetResultFromSerpApiByUsingKeywords(keywordSet, country, language)
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

// checkAndAuthInternal checks if the request includes internal info or not.
// If there is internal keys, validates them.
func checkAndAuthInternal(request events.APIGatewayProxyRequest) (bool, int, int, error) {
	accountName := ""
	if v, ok := request.QueryStringParameters["accountName"]; ok {
		accountName = v
	} else {
		return false, 0, http.StatusOK, nil
	}

	accountPassword := ""
	if v, ok := request.QueryStringParameters["accountPassword"]; ok {
		accountPassword = v
	} else {
		return false, 0, http.StatusUnauthorized, errors.New("Password is empty.")
	}

	internalJSON := os.Getenv("INTERNAL_ACCOUNTS_JSON")
	if internalJSON == "" {
		return false, 0, http.StatusUnauthorized, errors.New("Authorization is not valid.")
	}

	i := internal{}
	err := json.Unmarshal([]byte(internalJSON), &i)
	if err != nil {
		return false, 0, http.StatusInternalServerError, errors.New("We have some issues with excepting internal account logins. Please try later.")
	}

	h := sha256.New()
	h.Write([]byte(accountPassword))
	passwordHash := fmt.Sprintf("%x", h.Sum(nil))
	for _, x := range i.Accounts {
		if accountName == x.Name && passwordHash == x.Password {
			return true, x.Limit, http.StatusOK, nil
		}
	}

	return false, 0, http.StatusUnauthorized, errors.New("Authorization is not valid.")
}

// checkAndGetParams checks the params are set or not.
func checkAndSetParams(request events.APIGatewayProxyRequest) (int, error) {
	// Check the method.
	if request.HTTPMethod != "POST" {
		return http.StatusMethodNotAllowed, errors.New("Method not allowed. Only allowed: POST.")
	}

	rrType, err := getParam(request, "type")
	if err != nil {
		return http.StatusBadRequest, err
	}
	rType = rrType

	format, err = getParam(request, "format")
	if err != nil {
		return http.StatusBadRequest, err
	}

	country, err = getParam(request, "country")
	if err != nil {
		return http.StatusBadRequest, err
	}

	language, err = getParam(request, "language")
	if err != nil {
		return http.StatusBadRequest, err
	}

	return http.StatusOK, nil
}

// getParam returns params if it exists.
func getParam(request events.APIGatewayProxyRequest, param string) (string, error) {
	if v, ok := request.QueryStringParameters[param]; ok {
		return v, nil
	} else {
		return "", fmt.Errorf("%s is not set.", param)
	}
}

// checkLimit checks limit for the user.
func checkLimit(bodyLen int, isInternal bool, iLimit int) (int, error) {
	if !isInternal {
		// Check the count for non-login user.
		if bodyLen > 100 {
			return http.StatusBadRequest, errors.New("You have more than 100 URLs.")
		}
		return http.StatusOK, nil
	}

	switch iLimit {
	case -1: // limitless.
		return http.StatusOK, nil
	default: // has a limit.
		if bodyLen > iLimit {
			return http.StatusBadRequest, fmt.Errorf("You have more than %d URLs.", iLimit)
		}
	}

	return http.StatusOK, nil
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
