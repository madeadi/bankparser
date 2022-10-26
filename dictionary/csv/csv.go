package dictionary_csv

import (
	"encoding/csv"
	"log"
	"os"

	"github.com/madeadi/bankparser/internal"
)

type CsvStorage struct{}

func (d CsvStorage) Read(path string) map[string]string {
	// read the csv input
	arr := internal.ReadCsvFile(path)

	keywords := map[string]string{}
	for i, row := range arr {
		if i == 0 {
			// skip the header
			continue
		}
		keywords[row[0]] = row[1]
	}

	return keywords
}

func (d CsvStorage) Write(path string, data map[string]string) {
	csvFile, err := os.Create(path)
	if err != nil {
		panic(err)
	}
	defer csvFile.Close()

	csvWriter := csv.NewWriter(csvFile)
	defer csvWriter.Flush()

	// write the header
	if err := csvWriter.Write([]string{"Description", "Account"}); err != nil {
		panic(err)
	}

	for k, v := range data {
		if err := csvWriter.Write([]string{k, v}); err != nil {
			log.Fatalln("error writing record to file", err)
		}
	}
}
