package provider

import (
	"encoding/csv"
	"log"
	"os"
)

type Csv struct {
	Rows []Row
}

func (r Csv) getHeaders() []string {
	return []string{
		"Account",
		"Date",
		"Description",
		"In",
		"Out",
		"ContraAccount",
	}
}

func (c Csv) Save(outputCsv string) {
	output := outputCsv
	csvFile, err := os.Create(output)
	if err != nil {
		panic(err)
	}
	defer csvFile.Close()
	csvWriter := csv.NewWriter(csvFile)
	defer csvWriter.Flush()

	if err := csvWriter.Write(c.getHeaders()); err != nil {
		panic(err)
	}

	for _, record := range c.Rows {
		if err := csvWriter.Write(record.toSlice()); err != nil {
			log.Fatalln("error writing record to file", err)
		}
	}
}
