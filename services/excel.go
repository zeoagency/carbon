package services

import (
	"bytes"
	"fmt"

	"github.com/360EntSecGroup-Skylar/excelize/v2"
)

// ConvertURLResultToExcel creates a excel file by using success and fail variables.
func ConvertURLResultToExcel(success []URLOptionResponse, fail []string) (*bytes.Buffer, error) {
	f := excelize.NewFile()
	f.DeleteSheet("Sheet1") // delete default sheet.

	err := createSuccessSheetForURLs(f, success)
	if err != nil {
		return nil, err
	}
	err = createFailSheetForURLs(f, fail)
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
func createSuccessSheetForURLs(f *excelize.File, success []URLOptionResponse) error {
	f.NewSheet("success")

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
	for _, s := range success {
		err := f.SetCellValue("success", fmt.Sprintf("%s%d", letters[0], count), s.OriginalURL)
		if err != nil {
			return err
		}
		for i, url := range s.RelatedURLs {
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
func createFailSheetForURLs(f *excelize.File, fail []string) error {
	f.NewSheet("fail")

	err := f.SetCellValue("fail", "A1", "URL")
	if err != nil {
		return err
	}

	count := 2
	for _, fl := range fail {
		err := f.SetCellValue("fail", fmt.Sprintf("%s%d", "A", count), fl)
		if err != nil {
			return err
		}
		count++
	}

	return nil
}
