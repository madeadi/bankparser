package dictionary

import "testing"

type addTest struct {
	key, value, search string
}

var addTests = []addTest{
	{"That Coffee", "Coffee", "That Coffee"},
	{"Coffee Place", "Coffee", "Net*that Coffee Place"},
}

func TestDictionary(t *testing.T) {
	for _, v := range addTests {
		d := Dictionary{Keywords: map[string]string{}}
		d.AddKeyword(v.key, v.value)

		found := d.Translate(v.search)
		if found != v.value {
			t.Errorf("Searching: %s, Expected: %s, got: %s", v.search, v.value, found)
		}
	}
}
