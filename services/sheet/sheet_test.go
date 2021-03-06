package sheet

import (
	"fmt"
	"log"
	"testing"

	"github.com/joho/godotenv"

	"github.com/zeoagency/carbon/models"
	"github.com/zeoagency/carbon/services"
	"github.com/zeoagency/carbon/services/excel"
)

func init() {
	err := godotenv.Load("../../.env")
	if err != nil {
		log.Fatalln(err)
	}
}

func TestImportURLFileToGoogleSheets(t *testing.T) {
	urlSet := models.NewURLSet()
	urlSet.Add(
		"https://tools.zeo.org/carbon",
		"https://googlebunubulamaz.com/",
		"notaavalidurl",
		"https://boratanrikulu.dev/postgresql-nedir-nasil-calisir/",
		"https://boratanrikulu.dev/smtp-nasil-calisir-ve-postfix-kurulumu/",
	)

	_, err := services.GetResultByUsingURLs(urlSet, "tr", "tr")
	if err != nil {
		t.Fatal("Error occur while getting the result:", err)
	}

	f, err := excel.ConvertURLResultToExcel(urlSet)
	if err != nil {
		t.Fatal(err)
	}

	sheetURL, err := ImportFileToGoogleSheets(f)
	if err != nil {
		t.Fatalf("Test failed: %s", err)
	}

	fmt.Println(sheetURL)
}

func TestImportKeywordFileToGoogleSheets(t *testing.T) {
	keywordSet := models.NewKeywordSet()
	keywordSet.Add(
		"boratanrikulu blog archlinux install guide",
		"boratanrikulu blog postgresql nedir",
		"googlebunubulamaz blog",
	)

	_, err := services.GetResultByUsingKeywords(keywordSet, "tr", "tr")
	if err != nil {
		t.Fatal("Error occur while getting the result:", err)
	}

	f, err := excel.ConvertKeywordResultToExcel(keywordSet)
	if err != nil {
		t.Fatal(err)
	}

	sheetURL, err := ImportFileToGoogleSheets(f)
	if err != nil {
		t.Fatalf("Test failed: %s", err)
	}

	fmt.Println(sheetURL)
}
