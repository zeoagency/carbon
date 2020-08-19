package services

import (
	"fmt"
	"log"
	"testing"

	"github.com/joho/godotenv"

	"github.com/zeoagency/carbon/models"
)

func init() {
	err := godotenv.Load("../.env")
	if err != nil {
		log.Fatalln(err)
	}
}

func TestGetResultFromSerpApiByUsingURLs(t *testing.T) {
	urlSet := models.NewURLSet()

	urlSet.Add(
		"https://boratanrikulu.dev/postgresql-nedir-nasil-calisir/",
		"https://boratanrikulu.dev/smtp-nasil-calisir-ve-postfix-kurulumu/",
	)

	status, err := GetResultFromSerpApiByUsingURLs(urlSet, "tr", "tr")
	if err != nil {
		t.Fatalf("STATUS: %d ERROR: %s", status, err)
	}

	for originalURL, success := range urlSet.Successes {
		fmt.Printf("\n\tURL: %s\n", originalURL)
		for i, url := range success.URLs {
			fmt.Printf("\t\t%d - %s\n", i, url)
			fmt.Printf("\t\tSUGGESTED: %s\n", success.SuggestedURL)
		}
	}
}

func TestGetResultFromSerpApiByUsingKeywords(t *testing.T) {
	keywordSet := models.NewKeywordSet()

	keywordSet.Add(
		"boratanrikulu blog postgresql nedir nasil calisir",
		"boratanrikulu blog smtp nasil calisir",
	)

	status, err := GetResultFromSerpApiByUsingKeywords(keywordSet, "tr", "tr")
	if err != nil {
		t.Fatalf("STATUS: %d ERROR: %s", status, err)
	}

	for keyword, success := range keywordSet.Successes {
		fmt.Printf("\n\tURL: %s\n", keyword)
		for i, result := range success.Results {
			fmt.Printf("\t\t%d - %s\n", i, result)
		}
	}
}
