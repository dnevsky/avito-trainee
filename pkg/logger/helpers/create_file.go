package helpers

import (
	"encoding/csv"
	"io"
	"os"
)

const (
	csvPrefix = ".csv"
)

func GenerateTempFileName(format string) string {
	fileName, err := GenerateRandomString(10)
	if err != nil {
		return "temp" + format
	}
	return fileName + format
}

func DataToCsv(data [][]string) (file *os.File, err error) {
	fileName := GenerateTempFileName(csvPrefix)

	file, err = os.Create(fileName)
	if err != nil {
		return nil, err
	}

	writer := csv.NewWriter(file)

	if err = writer.WriteAll(data); err != nil {
		return nil, err
	}
	defer writer.Flush()

	_, err = file.Seek(0, io.SeekStart)
	return file, err
}
