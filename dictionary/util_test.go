package dictionary

import "testing"

func TestTokenise(t *testing.T) {
	input := "Net*that Coffee Place"
	expected := []string{
		"Net*that Coffee Place",
		"Net*that Coffee",
		"Coffee Place",
		"Net*that",
		"Coffee",
		"Place",
	}

	actual := tokenise(input)

	if len(actual) != len(expected) {
		t.Errorf("Expected %d tokens, got %d", len(expected), len(actual))
	}

	for i, v := range expected {
		if v != actual[i] {
			t.Errorf("Expected %s, got %s", v, actual[i])
		}
	}
}
