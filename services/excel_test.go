package services

import (
	"fmt"
	"testing"

	"github.com/360EntSecGroup-Skylar/excelize/v2"

	"github.com/zeoagency/carbon/models"
)

func TestConvertURLResultToExcel(t *testing.T) {
	urlSet := models.NewURLSet()
	urlSet.Add(
		"https://boratanrikulu.dev/archlinux-install",
		"https://boratanrikulu.dev/postgresql-nedir",
		"https://googlebunubulamaz.com/",
		"notaavalidurl",
	)
	_, err := GetResultByUsingURLs(urlSet, "tr", "tr")
	if err != nil {
		t.Fatal("Error occur while getting the result:", err)
	}

	f, err := ConvertURLResultToExcel(urlSet)
	if err != nil {
		t.Fatal(err)
	}

	eF, err := excelize.OpenReader(f)
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println("### SUCCESS SHEET ###")
	err = printSheet(eF, "success")
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println("### FAIL SHEET ###")
	err = printSheet(eF, "fail")
	if err != nil {
		t.Fatal(err)
	}
}

func TestConvertKeywordResultToExcel(t *testing.T) {
	keywordSet := models.NewKeywordSet()
	keywordSet.Add(
		"boratanrikulu blog archlinux install guide",
		"boratanrikulu blog postgresql nedir",
		"googlebunubulamaz blog",
	)
	_, err := GetResultByUsingKeywords(keywordSet, "tr", "tr")
	if err != nil {
		t.Fatal("Error occur while getting the result:", err)
	}

	f, err := ConvertKeywordResultToExcel(keywordSet)
	if err != nil {
		t.Fatal(err)
	}

	eF, err := excelize.OpenReader(f)
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println("### SUCCESS SHEET ###")
	err = printSheet(eF, "success")
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println("### FAIL SHEET ###")
	err = printSheet(eF, "fail")
	if err != nil {
		t.Fatal(err)
	}
}

// printSheet prints the sheet of the given excel.
func printSheet(f *excelize.File, sheet string) error {
	rows, err := f.GetRows(sheet)
	if err != nil {
		return err
	}

	for _, row := range rows {
		for _, colCell := range row {
			fmt.Print(colCell, "\t")
		}
		fmt.Println()
	}

	return nil
}
