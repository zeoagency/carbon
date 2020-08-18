package services

import (
	"bytes"
	"fmt"

	"github.com/360EntSecGroup-Skylar/excelize/v2"

	"gitlab.com/seo.do/zeo-carbon/helpers"
	"gitlab.com/seo.do/zeo-carbon/models"
)

var titleStyle = `{"alignment":{"horizontal":"center", "vertical":"center"}, "font":{"size":12, "name":"Calibri", "bold":true, "color":"#ffffff"}, "fill":{"type":"pattern", "color":["#000000"], "pattern":1}}`

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

	// Set styles
	err := f.SetColWidth("success", "A", "D", 40)
	if err != nil {
		return err
	}
	err = f.SetColWidth("success", "E", "E", 70)
	if err != nil {
		return err
	}
	style, err := f.NewStyle(titleStyle)
	if err != nil {
		return err
	}
	err = f.SetCellStyle("success", "A1", "E1", style)
	if err != nil {
		return err
	}

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
		for i, url := range success.URLs {
			err := f.SetCellValue("success", fmt.Sprintf("%s%d", letters[i+1], count), url)
			if err != nil {
				return err
			}
			err = f.SetCellValue("success", fmt.Sprintf("%s%d", letters[len(letters)-1], count), success.SuggestedURL)
			if err != nil {
				return err
			}
			index, _ := helpers.StringSliceContains(success.URLs, success.SuggestedURL)
			if index != -1 {
				style, err := f.NewStyle(`{"font":{"bold":true}}`)
				if err != nil {
					return err
				}
				f.SetCellStyle("success", fmt.Sprintf("%s%d", letters[index+1], count), fmt.Sprintf("%s%d", letters[index+1], count), style)
			}
		}
		count++
	}

	if count > 2 {
		style, err = f.NewStyle(`{"fill":{"type":"pattern", "color":["#d5eb81"], "pattern":1}}`)
		if err != nil {
			return err
		}
		err = f.SetCellStyle("success", "E2", fmt.Sprintf("E%d", count-1), style)
		if err != nil {
			return err
		}
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

	// Set styles
	err := f.SetColWidth("fail", "A", "B", 40)
	if err != nil {
		return err
	}
	style, err := f.NewStyle(titleStyle)
	if err != nil {
		return err
	}
	err = f.SetCellStyle("fail", "A1", "B1", style)
	if err != nil {
		return err
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
	letters := []string{"A", "B", "C", "D", "E"}
	titles := []string{"Keyword", "Position", "Title", "URL", "Description"}
	// NOTE: letters and titles sizes must be same!

	// Set titles.
	for i, letter := range letters {
		err := f.SetCellValue("success", fmt.Sprintf("%s%d", letter, 1), titles[i])
		if err != nil {
			return err
		}
	}

	// Set styles
	err := f.SetColWidth("success", "A", "A", 40)
	_ = f.SetColWidth("success", "C", "D", 40)
	_ = f.SetColWidth("success", "E", "E", 110)
	if err != nil {
		return err
	}
	style, err := f.NewStyle(titleStyle)
	if err != nil {
		return err
	}
	err = f.SetCellStyle("success", "A1", "E1", style)
	if err != nil {
		return err
	}

	count := 2
	for keyword, success := range keywordSet.Successes {
		err := f.SetCellValue("success", fmt.Sprintf("%s%d", letters[0], count), keyword)
		if err != nil {
			return err
		}
		groupCount := 0
		for i, result := range success.Results {
			if groupCount != 0 {
				location := fmt.Sprintf("%s%d", letters[0], count)
				err := f.SetCellValue("success", location, keyword)
				if err != nil {
					return err
				}
				style, err := f.NewStyle(`{"font":{"color":"#cccccc"}}`)
				if err != nil {
					return err
				}
				err = f.SetCellStyle("success", location, location, style)
				if err != nil {
					return err
				}
			}
			err = f.SetCellValue("success", fmt.Sprintf("%s%d", letters[1], count), fmt.Sprintf("#%d", i+1))
			if err != nil {
				return err
			}
			err = f.SetCellValue("success", fmt.Sprintf("%s%d", letters[2], count), result.Title)
			if err != nil {
				return err
			}
			err = f.SetCellValue("success", fmt.Sprintf("%s%d", letters[3], count), result.URL)
			if err != nil {
				return err
			}
			err = f.SetCellValue("success", fmt.Sprintf("%s%d", letters[4], count), result.Desc)
			if err != nil {
				return err
			}
			count++
			groupCount++
		}
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

	// Set styles
	err := f.SetColWidth("fail", "A", "B", 40)
	if err != nil {
		return err
	}
	style, err := f.NewStyle(titleStyle)
	if err != nil {
		return err
	}
	err = f.SetCellStyle("fail", "A1", "B1", style)
	if err != nil {
		return err
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
