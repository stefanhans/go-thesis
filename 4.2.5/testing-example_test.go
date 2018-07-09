package testingexample

import (
	"testing"
)

func Test_StringReverse(t *testing.T) {
	for _, testCase := range []struct {
		origin, expected string
	}{
		{"I am not reversed", "desrever ton ma I"},
		{"1234567890", "0987654321"},
		{"", ""},
	} {
		if Reverse(testCase.origin) != testCase.expected {
			t.Errorf("origin: %q => %q != expected: %q",
				testCase.origin,
				Reverse(testCase.origin),
				testCase.expected)
		}
	}
}
