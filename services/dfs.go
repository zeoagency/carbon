package services

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"sync"

	"github.com/zeoagency/carbon/helpers"
	"github.com/zeoagency/carbon/models"
)

type dfsApiResponse struct {
	StatusCode int `json:"status_code"`
	Tasks      []struct {
		Result []struct {
			Keyword string `json:"keyword"`
			Items   []struct {
				Type        string `json:"type"`
				Title       string `json:"title"`
				Description string `json:"description"`
				URL         string `json:"url"`
			} `json:"items"`
		} `json:"result"`
	} `json:"tasks"`
}

type dfsApiRequest struct {
	Keyword   string `json:"keyword"`
	Gl        string `json:"location_name"`
	Hl        string `json:"language_code"`
	SerpLimit int    `json:"depth"`
	Device    string `json:"device"`
}

// getResultFromDFSApi returns DFS API response for the given data.
func getResultFromDFSApi(kws keywords, country, language string, serpLimit int) ([]*dfsApiResponse, int, error) {
	responses := []*dfsApiResponse{}
	wg := new(sync.WaitGroup)

	for _, kw := range kws.ToStringSlice() {
		rq := []dfsApiRequest{}
		rq = append(rq, dfsApiRequest{
			Keyword:   kw,
			Gl:        "United Kingdom",
			Hl:        language,
			SerpLimit: serpLimit,
			Device:    "desktop",
		})

		wg.Add(1)
		go sendRequest(wg, &responses, rq)
	}

	wg.Wait()
	return responses, http.StatusCreated, nil
}

// sendRequest works as async, adds all results to the given slice.
func sendRequest(wg *sync.WaitGroup, responses *[]*dfsApiResponse, rq []dfsApiRequest) {
	defer wg.Done()

	rqJson, err := json.Marshal(rq)
	if err != nil {
		return
	}

	// Create the request.
	req, err := http.NewRequest("POST", "https://api.dataforseo.com/v3/serp/google/organic/live/regular", bytes.NewReader(rqJson))
	if err != nil {
		return
	}

	req.SetBasicAuth(os.Getenv("DFS_API_USER"), os.Getenv("DFS_API_PASSWORD"))

	// Send the request.
	c := &http.Client{}
	res, err := c.Do(req)
	if err != nil {
		return
	}
	defer res.Body.Close()

	// Check the result's status code.
	if !(res.StatusCode >= 200 && res.StatusCode <= 299) {
		log.Printf("Error: Unavailable SERP API Service. Status: %d\n", res.StatusCode)
		return
	}

	// Read the result, unmarshal it to the struct.
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return
	}

	response := dfsApiResponse{}
	err = json.Unmarshal(body, &response) // TODO: handle this error.

	*responses = append(*responses, &response)
}

// parseDFSResponseToFieldsForURLs extract the response to the URLSet.
// It only adds to success list when domains are matched.
// If it couldn't find any related URLs, it adds to the fail list.
func parseDFSResponseToFieldsForURLs(responses []*dfsApiResponse, urlSet *models.URLSet) {
	for _, response := range responses {
		for _, task := range response.Tasks {
			for _, result := range task.Result {
				url := urlSet.URLs[result.Keyword]

				r := []string{}
				for _, item := range result.Items {
					if len(r) == 3 {
						break
					}

					urlDomain, _, err := helpers.ExtractURL(item.URL)
					if err != nil {
						continue // The result's URL is not a valid URL.
					}
					if urlSet.URLs[result.Keyword].BaseURL == urlDomain {
						r = append(r, item.URL)
					}
				}

				if len(r) != 0 {
					urlSet.AddSuccess(url.FullURL, r)
				} else {
					urlSet.AddFail(url.FullURL, "We could not find any related URLs.")
				}
			}
		}
	}
}

// parseDFSResponseToFieldsForKeywords extract the response to the KeywordSet.
// It only adds to success list when the value is valid.
// If it couldn't find any results, it adds to the fail list.
func parseDFSResponseToFieldsForKeywords(responses []*dfsApiResponse, keywordSet *models.KeywordSet) {
	for _, response := range responses {
		for _, task := range response.Tasks {
			for _, result := range task.Result {
				r := []models.KeywordSuccessResult{}
				for _, item := range result.Items {
					if len(r) == 10 {
						break
					}
					r = append(r, models.KeywordSuccessResult{
						Title: item.Title,
						Desc:  item.Description,
						URL:   item.URL,
					})
				}
				if len(r) != 0 {
					keywordSet.AddSuccess(result.Keyword, r)
				} else {
					keywordSet.AddFail(result.Keyword, "We could not find any related URLs.")
				}
			}
		}
	}
}
