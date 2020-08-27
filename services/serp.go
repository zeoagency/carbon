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

	"github.com/zeoagency/carbon/helpers"
	"github.com/zeoagency/carbon/models"
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

// getResultFromSerpApi returns SERP API Response for the given data.
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

	api, selected := helpers.RandomAPICred()
	if selected == -1 {
		return nil, http.StatusServiceUnavailable, errors.New("We have some issues with the SERP API at this moment. Please try later.")
	}

	for i, _ := range api.Keys {
		if i >= 6 {
			break
		}
		// Example; If the selected is 4, Then it works like that: 4,5,6,7,8,9
		address, key := api.Keys[(selected+i)%len(api.Keys)].Address, api.Keys[(selected+i)%len(api.Keys)].Key

		// Create the request.
		req, err := http.NewRequest("POST", address, bytes.NewReader(rqJson))
		if err != nil {
			continue
		}

		// Set basic auth info.
		req.Header.Add("x-api-key", key)

		// Send the request.
		c := &http.Client{
			Timeout: 30 * time.Second,
		}
		res, err := c.Do(req)
		if err != nil {
			continue
		}
		defer res.Body.Close()

		// Check the result's status code.
		if !(res.StatusCode >= 200 && res.StatusCode <= 299) {
			log.Printf("Error: Unavailable SERP API Service. Status: %d\n", res.StatusCode)
			continue
		}

		// Read the result, unmarshal it to the struct.
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			continue
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

	}

	log.Println("Error: Unavailable SERP API Service.")
	return nil, http.StatusServiceUnavailable, errors.New("We have some issues with the SERP API at this moment. Please try later.")
}

// parseResponseToFieldsForURLs extract the response to the URLSet.
// It only adds to success list when domains are matched.
// If it couldn't find any related URLs, it adds to the fail list.
func parseSERPResponseToFieldsForURLs(response map[string][]serpApiResponse, urlSet *models.URLSet) {
	for _, url := range urlSet.URLs {
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
			urlSet.AddSuccess(url.FullURL, r)
			delete(urlSet.URLs, url.String())
		} else {
			urlSet.AddFail(url.FullURL, "We could not find any related URLs.")
		}
	}
}

// parseResponseToFieldsForKeywords extract the response to the KeywordSet.
// It only adds to success list when the value is valid.
// If it couldn't find any results, it adds to the fail list.
func parseSERPResponseToFieldsForKeywords(response map[string][]serpApiResponse, keywordSet *models.KeywordSet) {
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
			delete(keywordSet.Keywords, keyword)
		} else {
			keywordSet.AddFail(keyword, "We could not find any result.")
		}
	}
}
