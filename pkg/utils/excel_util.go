package utils

import (
	"mime/multipart"

	"github.com/xuri/excelize/v2"
)

func ExcelExtract(file *multipart.FileHeader) (map[string][]string, error) {
	src, err := file.Open()
	if err != nil {
		return nil, err
	}
	defer src.Close()

	xlsx, err := excelize.OpenReader(src)
	if err != nil {
		return nil, err
	}

	sheets := xlsx.GetSheetList()
	firstSheetName := sheets[0]

	rows, err := xlsx.GetRows(firstSheetName)
	if err != nil {
		return nil, err
	}

	result := make(map[string][]string)
	if len(rows) > 1 {
		headers := rows[0]

		for i, header := range headers {
			data := []string{}

			for _, row := range rows[1:] {
				if i < len(row) {
					data = append(data, row[i])
				}
			}
			result[header] = data
		}
	}

	return result, nil
}
