package bar

import (
	"reflect"
	"testing"
)

func TestTokenize(t *testing.T) {
	var testCases = []struct {
		formatString string
		expected     tokens
	}{
		{"", nil},
		{" ", tokens{spaceToken{}}},
		{":bar", tokens{barToken{}}},
		{" :bar", tokens{spaceToken{}, barToken{}}},
		{" :bar ", tokens{spaceToken{}, barToken{}, spaceToken{}}},
		{"  :bar", tokens{spaceToken{}, spaceToken{}, barToken{}}},
		{":bar:bar", tokens{barToken{}, barToken{}}},
		{":bar:rate", tokens{barToken{}, rateToken{}}},
		{"bar", tokens{literalToken{"bar"}}},
		{"bar:bar", tokens{literalToken{"bar"}, barToken{}}},
		{"不与", tokens{literalToken{"不与"}}},
		{"不与:bar", tokens{literalToken{"不与"}, barToken{}}},
	}

	for i, testCase := range testCases {
		got := tokenize(testCase.formatString, nil)
		if !reflect.DeepEqual(got, testCase.expected) {
			t.Errorf(
				"[%d] tokenize(%#v, nil)\n\n  got %#v\n  want %#v",
				i,
				testCase.formatString,
				got,
				testCase.expected,
			)
		}
	}
}

func TestTokenizeWithBoundaryCharacters(t *testing.T) {
	var testCases = []struct {
		formatString string
		expected     tokens
	}{
		{"(:bar", tokens{literalToken{"("}, barToken{}}},
		{"(:bar)", tokens{literalToken{"("}, barToken{}, literalToken{")"}}},
		{":bar (:eta remaining)", tokens{barToken{}, spaceToken{}, literalToken{"("}, etaToken{}, spaceToken{}, literalToken{"remaining)"}}},
	}

	for i, testCase := range testCases {
		got := tokenize(testCase.formatString, nil)
		if !reflect.DeepEqual(got, testCase.expected) {
			t.Errorf(
				"[%d] tokenize(%#v, nil)\n\n  got %#v\n  want %#v",
				i,
				testCase.formatString,
				got,
				testCase.expected,
			)
		}
	}
}

func TestTokenizeWithCustomVerbs(t *testing.T) {
	var testCases = []struct {
		formatString string
		customVerbs  []string
		expected     tokens
	}{
		{":custom", nil, tokens{literalToken{":custom"}}},
		{":custom", []string{"custom"}, tokens{customVerbToken{"custom"}}},
		{":bar:custom", []string{"custom"}, tokens{barToken{}, customVerbToken{"custom"}}},
	}

	for i, testCase := range testCases {
		got := tokenize(testCase.formatString, testCase.customVerbs)
		if !reflect.DeepEqual(got, testCase.expected) {
			t.Errorf(
				"[%d] tokenize(%#v, %#v)\n\n  got %#v\n  want %#v",
				i,
				testCase.formatString,
				testCase.customVerbs,
				got,
				testCase.expected,
			)
		}
	}
}
