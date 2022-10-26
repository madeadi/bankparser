package main

import (
	"flag"
	"fmt"
	"os"

	_dict "github.com/madeadi/bankparser/dictionary"
	_csv "github.com/madeadi/bankparser/dictionary/csv"
	"github.com/madeadi/bankparser/provider"
	"github.com/madeadi/bankparser/provider/revolut"
)

func main() {
	providerStr := flag.String("provider", "revolut", "The provider of the CSV file")
	in := flag.String("in", "", "The input CSV file")
	out := flag.String("out", "", "The output CSV file")
	dict := flag.String("dict", "", "The output CSV file")
	flag.Parse()

	if *in == "" {
		panic("No input path specified")
	}
	if *out == "" {
		panic("No output path specified")
	}

	dictionary := _dict.Dictionary{Keywords: make(map[string]string)}
	if *dict != "" {
		dictionaryPath := "./dictionary.csv"
		dictStorage := _csv.CsvStorage{}
		dictionary = _dict.NewDictionary(dictionaryPath, dictStorage)
	}

	var p provider.Provider

	switch *providerStr {
	case "revolut":
		p = revolut.NewRevolut(dictionary)
		p.Parse(*in).Save(*out)

	default:
		panic("Unknown provider")
	}

	fmt.Println("Successfully parsing to " + *out)
}

/// Save keywords to keywords file
func saveKeywords(inputPath *string, outputPath *string, p provider.Provider) {
	keywords := p.GenerateKeywords(*inputPath)
	kf, err := os.OpenFile(*outputPath, os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}
	defer kf.Close()

	for _, v := range keywords {
		kf.WriteString(fmt.Sprintf("%s\n", v))
	}
}
