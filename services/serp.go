package services

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"os"
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
func GetResultFromSerpApiByUsingKeywords(keywords *models.KeywordSet, language string) (int, error) {
	response, status, err := getResultFromSerpApi(keywords, language, 20)
	if err != nil {
		return status, err
	}

	parseResponseToFieldsForKeywords(response, keywords)
	return http.StatusOK, nil
}

// GetResultFromSerpApiByUsingURLs add the result to the given URLSet by talking the Serp API.
func GetResultFromSerpApiByUsingURLs(urls *models.URLSet, language string) (int, error) {
	response, status, err := getResultFromSerpApi(urls, language, 10)
	if err != nil {
		return status, err
	}

	parseResponseToFieldsForURLs(response, urls)
	return http.StatusOK, nil
}

// getResultFromSerpApi returns SERP API Response for given data type.
func getResultFromSerpApi(kws keywords, language string, serpLimit int) (map[string][]serpApiResponse, int, error) {
	// Create the request body.
	rq := serpApiRequest{
		Keywords:  kws.ToStringSlice(),
		Gl:        language,
		Hl:        language,
		SerpLimit: strconv.Itoa(serpLimit),
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
	c := &http.Client{
		Timeout: 30 * time.Second,
	}
	res, err := c.Do(req)
	defer res.Body.Close()
	if err != nil {
		return nil, http.StatusServiceUnavailable, errors.New("Error: We have some issues with SERP API at this moment. Please try later.")
	}

	if !(res.StatusCode >= 200 && res.StatusCode <= 299) {
		log.Printf("Error: Unavailable SERP API Service.\nStatus: %d\n", res.StatusCode)
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
