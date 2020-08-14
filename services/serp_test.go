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

	result, status, err := GetResultFromSerpApiByUsingURLs(urls, "tr")
	if err != nil {
		t.Fatalf("STATUS: %d ERROR: %s", status, err)
	}

	for key, value := range result {
		fmt.Printf("\n\tKEYWORD: %s\n\tURL:\n", key)
		for i, url := range value {
			fmt.Printf("\t\t%d - %s\n", i, url)
		}
	}
}
