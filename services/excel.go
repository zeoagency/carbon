package services

import (
	"bytes"
	"fmt"

	"github.com/360EntSecGroup-Skylar/excelize/v2"

	"gitlab.com/seo.do/zeo-carbon/models"
)

// ConvertURLResultToExcel creates a excel file by using the URLSet.
func ConvertURLResultToExcel(urlSet *models.URLSet) (*bytes.Buffer, error) {
	f := excelize.NewFile()
	f.NewSheet("success")
	f.NewSheet("fail")
	f.DeleteSheet("Sheet1") // delete default sheet.

	err := createSuccessSheetForURLs(f, urlSet)
	if err != nil {
		return nil, err
	}

	err = createFailSheetForURLs(f, urlSet)
	if err != nil {
		return nil, err
	}

	b, err := f.WriteToBuffer()
	if err != nil {
		return nil, err
	}

	return b, nil
}

// ConvertKeywordResultToExcel creates a excel file by using the KeywordSet.
func ConvertKeywordResultToExcel(keywordSet *models.KeywordSet) (*bytes.Buffer, error) {
	f := excelize.NewFile()
	f.NewSheet("success")
	f.NewSheet("fail")
	f.DeleteSheet("Sheet1") // delete default sheet.

	err := createSuccessSheetForKeywords(f, keywordSet)
	if err != nil {
		return nil, err
	}

	err = createFailSheetForKeywords(f, keywordSet)
	if err != nil {
		return nil, err
	}

	b, err := f.WriteToBuffer()
	if err != nil {
		return nil, err
	}

	return b, nil
}

// createSuccessSheetForURLs creates success sheet for the given excel.
func createSuccessSheetForURLs(f *excelize.File, urlSet *models.URLSet) error {
	letters := []string{"A", "B", "C", "D", "E"}
	titles := []string{"URL", "Alternative 1", "Alternative 2", "Alternative 3", "Suggested"}
	// NOTE: letters and titles sizes must be same!

	// Set titles.
	for i, letter := range letters {
		err := f.SetCellValue("success", fmt.Sprintf("%s%d", letter, 1), titles[i])
		if err != nil {
			return err
		}
	}

	count := 2
	for originalURL, success := range urlSet.Successes {
		err := f.SetCellValue("success", fmt.Sprintf("%s%d", letters[0], count), originalURL)
		if err != nil {
			return err
		}
		for i, url := range success.RelatedURLs {
			err := f.SetCellValue("success", fmt.Sprintf("%s%d", letters[i+1], count), url)
			if err != nil {
				return err
			}
		}
		count++
	}

	return nil
}

// createFailSheetForURLs creates fail sheet for the given excel.
func createFailSheetForURLs(f *excelize.File, urlSet *models.URLSet) error {
	letters := []string{"A", "B"}
	titles := []string{"URL", "Reason"}
	// NOTE: letters and titles sizes must be same!

	// Set titles.
	for i, letter := range letters {
		err := f.SetCellValue("fail", fmt.Sprintf("%s%d", letter, 1), titles[i])
		if err != nil {
			return err
		}
	}

	count := 2
	for originalURL, fail := range urlSet.Fails {
		err := f.SetCellValue("fail", fmt.Sprintf("%s%d", "A", count), originalURL)
		if err != nil {
			return err
		}
		err = f.SetCellValue("fail", fmt.Sprintf("%s%d", "B", count), fail.Reason)
		if err != nil {
			return err
		}
		count++
	}

	return nil
}

// createSuccessSheetForKeywords creates success sheet for the given excel.
func createSuccessSheetForKeywords(f *excelize.File, keywordSet *models.KeywordSet) error {
	letters := []string{"A", "B", "C", "D", "E", "F", "G", "H", "J", "K", "L"}
	titles := []string{"URL", "# 1", "# 2", "# 3", "# 4", "# 5", "# 6", "# 7", "# 8", "# 9", "# 10"}
	// NOTE: letters and titles sizes must be same!

	// Set titles.
	for i, letter := range letters {
		err := f.SetCellValue("success", fmt.Sprintf("%s%d", letter, 1), titles[i])
		if err != nil {
			return err
		}
	}

	count := 2
	for keyword, success := range keywordSet.Successes {
		err := f.SetCellValue("success", fmt.Sprintf("%s%d", letters[0], count), keyword)
		if err != nil {
			return err
		}
		for i, result := range success.Results {
			err := f.SetCellValue("success", fmt.Sprintf("%s%d", letters[i+1], count),
				fmt.Sprintf("Title: %s - Desc: %s - URL: %s",
					result.Title, result.Desc, result.URL,
				))
			if err != nil {
				return err
			}
		}
		count++
	}

	return nil
}

// createFailSheetForKeywords creates fail sheet for the given excel.
func createFailSheetForKeywords(f *excelize.File, keywordSet *models.KeywordSet) error {
	letters := []string{"A", "B"}
	titles := []string{"Keyword", "Reason"}
	// NOTE: letters and titles sizes must be same!

	// Set titles.
	for i, letter := range letters {
		err := f.SetCellValue("fail", fmt.Sprintf("%s%d", letter, 1), titles[i])
		if err != nil {
			return err
		}
	}

	count := 2
	for keyword, fail := range keywordSet.Fails {
		err := f.SetCellValue("fail", fmt.Sprintf("%s%d", "A", count), keyword)
		if err != nil {
			return err
		}
		err = f.SetCellValue("fail", fmt.Sprintf("%s%d", "B", count), fail.Reason)
		if err != nil {
			return err
		}
		count++
	}

	return nil
}
