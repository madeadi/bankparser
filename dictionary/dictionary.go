package dictionary

import "strings"

type Dictionary struct {
	Keywords map[string]string
}

func NewDictionary(path string, storage Storage) Dictionary {
	dict := Dictionary{}
	dict.Keywords = make(map[string]string)

	dict.Keywords = storage.Read(path)

	for k, v := range dict.Keywords {
		dict.Keywords[strings.ToLower(k)] = v
	}

	return dict
}

func (d Dictionary) AddKeywords(keywords map[string]string) {
	for k, v := range keywords {
		d.AddKeyword(k, v)
	}
}

func (d Dictionary) AddKeyword(keyword string, account string) {
	if _, ok := d.Keywords[keyword]; ok {
		return
	}
	d.Keywords[keyword] = account
}

func (d Dictionary) Translate(input string) string {
	tokens := tokenise(input)
	for _, v := range tokens {
		if val, ok := d.Keywords[v]; ok {
			return val
		}
	}

	return ""
}
