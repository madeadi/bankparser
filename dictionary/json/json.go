package json

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

type entry struct {
	Keyword string `json:"keyword"`
	Account string `json:"account"`
}
type JsonStorage struct {
}

func (d JsonStorage) Read(path string) map[string]string {
	jsonFile, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer jsonFile.Close()

	byteValue, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		panic(err)
	}

	entries := []entry{}
	if err := json.Unmarshal(byteValue, &entries); err != nil {
		panic(err)
	}

	keywords := map[string]string{}
	for _, e := range entries {
		keywords[e.Keyword] = e.Account
	}

	return keywords
}

func (d JsonStorage) Write(path string, data map[string]string) {
	entries := []entry{}
	for k, v := range data {
		entries = append(entries, entry{
			Keyword: k,
			Account: v,
		})
	}

	file, err := json.MarshalIndent(entries, "", " ")
	if err != nil {
		panic(err)
	}

	err = ioutil.WriteFile(path, file, 0644)
	if err != nil {
		panic(err)
	}
}
