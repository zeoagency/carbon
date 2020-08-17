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

func TestImportFileToGoogleSheets(t *testing.T) {
	urlSet := models.NewURLSet()
	urlSet.Add(
		"https://boratanrikulu.dev/archlinux-install",
		"https://boratanrikulu.dev/postgresql-nedir",
		"https://googlebunubulamaz.com/",
		"notaavalidurl",
	)

	_, err := GetResultFromSerpApiByUsingURLs(urlSet, "tr")
	if err != nil {
		t.Fatal("Error occur while getting the result:", err)
	}

	f, err := ConvertURLResultToExcel(urlSet)
	if err != nil {
		t.Fatal(err)
	}

	sheetURL, err := ImportFileToGoogleSheets(f)
	if err != nil {
		t.Fatalf("Test failed: %s", err)
	}

	fmt.Println(sheetURL)
}
