package revolut

import (
	"strconv"

	"github.com/madeadi/bankparser/dictionary"
	"github.com/madeadi/bankparser/internal"
	"github.com/madeadi/bankparser/provider"
)

type Revolut struct {
	outputPath string
	dictionary dictionary.Dictionary
}

func NewRevolut(d dictionary.Dictionary) provider.Provider {
	return Revolut{
		dictionary: d,
	}
}

func (r Revolut) Parse(inputCsv string) *provider.Csv {
	input := inputCsv
	res := internal.ReadCsvFile(input)

	var txns []provider.Row
	for i, row := range res {
		// ignore header
		if i == 0 {
			continue
		}
		var (
			in  float64
			out float64
		)
		amount, _ := strconv.ParseFloat(row[5], 64)
		fee, _ := strconv.ParseFloat(row[6], 64)
		if amount < 0 {
			out = -amount + fee
		} else {
			in = amount
		}

		description := row[4]
		txn := provider.Row{
			Account:       "revolut",
			Date:          row[2],
			Description:   description,
			In:            in,
			Out:           out,
			ContraAccount: r.dictionary.Translate(description),
		}

		txns = append(txns, txn)
	}

	return &provider.Csv{
		Rows: txns,
	}
}

func (r Revolut) GenerateKeywords(inputCsv string) []string {
	csv := r.Parse(inputCsv)

	maps := map[string]int{}
	for _, row := range csv.Rows {
		maps[row.Description] = 1
	}

	keywords := make([]string, len(maps))
	i := 0
	for k := range maps {
		keywords[i] = k
		i++
	}

	return keywords
}
