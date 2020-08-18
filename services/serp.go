package services

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"

	"gitlab.com/seo.do/zeo-carbon/helpers"
	"gitlab.com/seo.do/zeo-carbon/models"
)

// The API response includes this struct as an array for each keywords.
// So, the real response like this: `map[string][]serpApiResponse{}`
type serpApiResponse struct {
	Result struct {
		Left []struct {
			Title   string `json:"title,omitempty"`
			Type    string `json:"type"`
			URL     string `json:"url,omitempty"`
			Snippet string `json:"snippet"`
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
func GetResultFromSerpApiByUsingKeywords(keywords *models.KeywordSet, country, language string) (int, error) {
	response, status, err := getResultFromSerpApi(keywords, country, language, 20)
	if err != nil {
		return status, err
	}

	parseResponseToFieldsForKeywords(response, keywords)
	return http.StatusOK, nil
}

// GetResultFromSerpApiByUsingURLs add the result to the given URLSet by talking the Serp API.
func GetResultFromSerpApiByUsingURLs(urls *models.URLSet, country, language string) (int, error) {
	response, status, err := getResultFromSerpApi(urls, country, language, 10)
	if err != nil {
		return status, err
	}

	parseResponseToFieldsForURLs(response, urls)
	return http.StatusOK, nil
}

// getResultFromSerpApi returns SERP API Response for given data type.
func getResultFromSerpApi(kws keywords, country, language string, serpLimit int) (map[string][]serpApiResponse, int, error) {
	// Create the request body.
	rq := serpApiRequest{
		Keywords:  kws.ToStringSlice(),
		Gl:        country,
		Hl:        language,
		SerpLimit: strconv.Itoa(serpLimit),
		Device:    "desktop",
	}

	rqJson, err := json.Marshal(rq)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	// Let's start to sending requests.
	// It will try 10 times at most.
	// For each time, randomly API address will be selected.
	tries := 0
	for tries < 10 {
		// Get API address and key..
		address, key, err := helpers.RandomAPICred()
		if err != nil {
			return nil, http.StatusInternalServerError, errors.New("We have some issues with the SERP API at this moment. Please try later.")
		}

		// Create the request.
		req, err := http.NewRequest("POST", address, bytes.NewReader(rqJson))
		if err != nil {
			return nil, http.StatusInternalServerError, errors.New("We have some issues with the SERP API at this moment. Please try later.")
		}

		// Set basic auth info.
		req.Header.Add("x-api-key", key)

		// Send the request.
		c := &http.Client{
			Timeout: 30 * time.Second,
		}
		res, err := c.Do(req)
		defer res.Body.Close()
		if err != nil {
			return nil, http.StatusServiceUnavailable, errors.New("We have some issues with SERP API at this moment. Please try later.")
		}

		// Check the result's status code.
		if !(res.StatusCode >= 200 && res.StatusCode <= 299) {
			log.Printf("Error: Unavailable SERP API Service.\nStatus: %d\n", res.StatusCode)
			return nil, http.StatusServiceUnavailable, errors.New("We have some issues with SERP API at this moment. Please try later.")
		}

		// Read the result, unmarshal it to the struct.
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return nil, http.StatusInternalServerError, errors.New("We have some issues at this moment. Please try later.")
		}
		rsMap := make(map[string]map[string][]serpApiResponse)
		_ = json.Unmarshal(body, &rsMap) // TODO: handle this error.

		// Check if there is any result?!.
		// If there is no result, we will be on the roads one more time.
		for _, v := range rsMap {
			for _, k := range v {
				if len(k) != 0 {
					return rsMap["data"], http.StatusOK, nil
				}
			}

		}

		tries++
	}

	return nil, http.StatusServiceUnavailable, errors.New("We have some issues with the SERP API at this moment. Please try later.")
}

// parseResponseToFieldsForURLs extract the response to the URLSet.
// It only adds to success list when domains are matched.
// If it couldn't find any related URLs, it adds to the fail list.
func parseResponseToFieldsForURLs(response map[string][]serpApiResponse, urlSet *models.URLSet) {
	for _, url := range urlSet.URLs {
		originalURL := url.FullURL
		r := []string{}
		for _, value := range response[url.String()] {
			for _, v := range value.Result.Left {
				// Stop adding if there is already 3 URLs.
				if len(r) == 3 {
					break
				}

				urlDomain, _, err := helpers.ExtractURL(v.URL)
				if err != nil {
					continue // The result's URL is not a valid URL.
				}
				if v.Type == "organic" && url.BaseURL == urlDomain {
					r = append(r, v.URL)
				}
			}
		}
		// Add results to the lists.
		if len(r) != 0 {
			urlSet.AddSuccess(originalURL, r)
		} else {
			urlSet.AddFail(originalURL, "We could not find any related URLs.")
		}
	}
}

// parseResponseToFieldsForKeywords extract the response to the KeywordSet.
// It only adds to success list when the value is valid.
// If it couldn't find any results, it adds to the fail list.
func parseResponseToFieldsForKeywords(response map[string][]serpApiResponse, keywordSet *models.KeywordSet) {
	for keyword, _ := range keywordSet.Keywords {
		r := []models.KeywordSuccessResult{}
		for _, value := range response[keyword] {
			for _, v := range value.Result.Left {
				// Stop adding if there is already 10 results.
				if len(r) == 10 {
					break
				}

				if v.Type == "organic" && v.URL != "" {
					r = append(r, models.KeywordSuccessResult{
						Title: v.Title,
						Desc:  v.Snippet,
						URL:   v.URL,
					})
				}
			}
		}
		// Add results to the lists.
		if len(r) != 0 {
			keywordSet.AddSuccess(keyword, r)
		} else {
			keywordSet.AddFail(keyword, "We could not find any result.")
		}
	}
}
