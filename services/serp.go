package services

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"

	"gitlab.com/seo.do/zeo-carbon/models"
)

// The API response includes this struct as an array for each keywords.
// So, the real response like this: `map[string][]serpApiResponse{}`
type serpApiResponse struct {
	Result struct {
		Left []struct {
			Rank  int    `json:"rank"`
			Title string `json:"title,omitempty"`
			Type  string `json:"type"`
			URL   string `json:"url,omitempty"`
		} `json:"left"`
	} `json:"result"`
}

type serpApiRequest struct {
	Keywords  []string `json:"keyword"`
	Gl        string   `json:"gl"`
	Hl        string   `json:"hl"`
	SerpLimit string   `json:"serp_limit"`
	Device    string   `json:"device"`
}

// keywords is an interface that includes ToStringSlice method.
// You can use models.URLset or models.KeywordSet for this interface.
type keywords interface {
	ToStringSlice() []string
}

// GetResultFromSerpApiByUsingKeywords returns related 10 results for each Keywords by talking with SEPR API.
func GetResultFromSerpApiByUsingKeywords(keywords *models.KeywordSet, language string) {
	// TODO...
}

// GetResultFromSerpApiByUsingURLs returns related 3 results for each URLs by talking with SERP API.
// The result value is a map that's contains 3 URLs for each URLs.
// It returns a status code and error message if there is an any issue.
func GetResultFromSerpApiByUsingURLs(urls *models.URLSet, language string) (map[string][]string, int, error) {
	response, status, err := getResultFromSerpApi(urls, language)
	if err != nil {
		return nil, status, err
	}

	result := parseResponseToMapForURLs(response, urls)
	return result, http.StatusOK, nil
}

// getResultFromSerpApi returns SERP Api Response for given data type.
func getResultFromSerpApi(kws keywords, language string) (map[string][]serpApiResponse, int, error) {
	fmt.Println(kws)

	// Create the request body.
	rq := serpApiRequest{
		Keywords:  kws.ToStringSlice(),
		Gl:        language,
		Hl:        language,
		SerpLimit: "10",
		Device:    "desktop",
	}

	rqJson, err := json.Marshal(rq)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	// Create the request.
	req, err := http.NewRequest(
		"POST", os.Getenv("SERP_API_ADDRESS"),
		bytes.NewReader(rqJson),
	)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	// Set basic auth info.
	req.SetBasicAuth(os.Getenv("SERP_API_USERNAME"), os.Getenv("SERP_API_PASSWORD"))

	// Send the request.
	c := &http.Client{}
	res, err := c.Do(req)
	if err != nil {
		return nil, http.StatusServiceUnavailable, errors.New("Error: We have some issues while sending request.")
	}
	defer res.Body.Close()

	if !(res.StatusCode >= 200 && res.StatusCode <= 299) {
		log.Printf("Error: Unavailable SERP API Service. Status: %d", res.StatusCode)
		return nil, http.StatusServiceUnavailable, errors.New("Error: We have some issues with SERP API at this moment. Please try later.")
	}

	// Read the result, unmarshal it to the struct.
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}
	rs := map[string][]serpApiResponse{}
	err = json.Unmarshal(body, &rs)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	return rs, http.StatusOK, nil
}

// parseResponseToMapForURLs converts map[string][]serpApiResponse to map[string][]string.
// It selects related -organic- 3 URLs for each URLs.
// It is used to parse response for URLs option.
func parseResponseToMapForURLs(response map[string][]serpApiResponse, urlSet *models.URLSet) map[string][]string {
	result := map[string][]string{}

	for key, value := range response {
		for _, rs := range value {
			r := []string{}

			for _, v := range rs.Result.Left {
				// Stop adding if there is already 3 URLs.
				if len(r) == 3 {
					break
				}

				urlParsed, err := url.Parse(v.URL)
				if err != nil {
					continue
				}
				if v.Type == "organic" && urlSet.URLs[key].BaseURL == urlParsed.Hostname() {
					r = append(r, v.URL)
				}
			}

			result[key] = r
		}
	}

	return result
}

// parseResponseToMapForURLs converts map[string][]serpApiResponse to map[string]interface{}.
// It selects related -organic- 10 URLs for each keywords.
// It is used to parse response for Keywords option.
// The interface may include: Title, Desc and URL.
func parseResponseToMapForKeywords(response map[string][]serpApiResponse, urlSet *models.URLSet) map[string]interface{} {
	// TODO...
	return nil
}
