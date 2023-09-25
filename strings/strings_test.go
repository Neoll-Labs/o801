/*
 license x
*/

package strings

import "testing"

func TestPlural(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		input    string
		expected string
	}{
		{"apple", "apples"},
		{"car", "cars"},
		{"dog", "dogs"},
		{"", "s"},
	}

	for _, tp := range testCases {
		tc := tp
		t.Run(tc.input, func(t *testing.T) {
			result := Plural(tc.input)
			if result != tc.expected {
				t.Errorf("Expected: %s, Got: %s", tc.expected, result)
			}
		})
	}
}
