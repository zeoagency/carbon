package services

import (
	"net/http"

	"github.com/zeoagency/carbon/models"
)

// keywords is an interface that includes ToStringSlice method.
// You can use models.URLset or models.KeywordSet for this interface.
// It used in SERP and DFS.
type keywords interface {
	ToStringSlice() []string
}

// GetResultByUsingURLs add the result to the given URLSet by talking with Serp API or DFS.
func GetResultByUsingURLs(urls *models.URLSet, country, language string) (int, error) {
	response, _, err := getResultFromSerpApi(urls, country, language, 10)
	if err == nil {
		parseSERPResponseToFieldsForURLs(response, urls)
		if len(urls.URLs) == 0 {
			return http.StatusOK, nil
		}
	}

	// That means there is still unprocessed URLs exist.

	dfsresponse, status, err := getResultFromDFSApi(urls, country, language, 10)
	if err != nil {
		return status, err
	}

	parseDFSResponseToFieldsForURLs(dfsresponse, urls)
	return http.StatusOK, nil
}

// GetResultByUsingKeywords returns related 10 results for each Keywords by talking with SEPR API or DFS.
func GetResultByUsingKeywords(keywords *models.KeywordSet, country, language string) (int, error) {
	response, _, err := getResultFromSerpApi(keywords, country, language, 10)
	if err == nil {
		parseSERPResponseToFieldsForKeywords(response, keywords)
		if len(keywords.Keywords) == 0 {
			return http.StatusOK, nil
		}
	}

	// That means there is still unprocessed Keywords exist.

	dfsresponse, status, err := getResultFromDFSApi(keywords, country, language, 10)
	if err != nil {
		return status, err
	}

	parseDFSResponseToFieldsForKeywords(dfsresponse, keywords)
	return http.StatusOK, nil
}
