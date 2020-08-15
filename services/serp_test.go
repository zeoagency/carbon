package services

import (
	"fmt"
	"log"
	"testing"

	"github.com/joho/godotenv"

	"gitlab.com/seo.do/zeo-carbon/models"
)

func init() {
	err := godotenv.Load("../.env")
	if err != nil {
		log.Fatalln(err)
	}
}

func TestGetResultFromSerpApiByUsingURLs(t *testing.T) {
	urls := models.NewURLSet()

	err := urls.Add(
		"https://boratanrikulu.dev/postgresql-nedir-nasil-calisir/",
		"https://boratanrikulu.dev/smtp-nasil-calisir-ve-postfix-kurulumu/",
	)
	if err != nil {
		t.Fatal("Error occur while creating URLSet.")
	}

	response, _, status, err := GetResultFromSerpApiByUsingURLs(urls, "tr")
	if err != nil {
		t.Fatalf("STATUS: %d ERROR: %s", status, err)
	}

	for _, r := range response {
		fmt.Printf("\n\tURL: %s\n", r.OriginalURL)
		for i, url := range r.RelatedURLs {
			fmt.Printf("\t\t%d - %s\n", i, url)
		}
	}
}
