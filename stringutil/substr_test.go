package stringutil

import "testing"

func TestSubstrBeforeHandler(t *testing.T) {
	var tests = []struct {
		str  string
		find string
		want string
	}{
		{"", "abc", ""},
		{"abc", "a", ""},
		{"abc", "b", "a"},
		{"abc", "c", "ab"},
		{"abc", "d", "abc"},
		{"abc", "", ""},
	}

	for _, test := range tests {
		if got := SubstrBefore(test.str, test.find); got != test.want {
			t.Errorf("SubstrBefore(%q, %q) = %q; want: %q", test.str, test.find, got, test.want)
		}
	}
}

func TestSubstrAfterHandler(t *testing.T) {
	var tests = []struct {
		str  string
		find string
		want string
	}{
		{"", "*", ""},
		{"abc", "a", "bc"},
		{"abcba", "b", "cba"},
		{"abc", "c", ""},
		{"abc", "d", ""},
		{"abc", "", "abc"},
	}

	for _, test := range tests {
		if got := SubstrAfter(test.str, test.find); got != test.want {
			t.Errorf("SubstrAfter(%q, %q) = %q; want: %q", test.str, test.find, got, test.want)
		}
	}
}

func TestSubstrBeforeLastHandler(t *testing.T) {
	var tests = []struct {
		str  string
		find string
		want string
	}{
		{"", "*", ""},
		{"abcba", "b", "abc"},
		{"abc", "c", "ab"},
		{"a", "a", ""},
		{"a", "z", "a"},
		{"a", "", "a"},
	}

	for _, test := range tests {
		if got := SubstrBeforeLast(test.str, test.find); got != test.want {
			t.Errorf("SubstrBeforeLast(%q, %q) = %q; want: %q", test.str, test.find, got, test.want)
		}
	}
}

func TestSubstrAfterLastHandler(t *testing.T) {
	var tests = []struct {
		str  string
		find string
		want string
	}{
		{"", "*", ""},
		{"*", "", ""},
		{"abc", "a", "bc"},
		{"abcba", "b", "a"},
		{"abc", "c", ""},
		{"a", "a", ""},
		{"a", "z", ""},
	}

	for _, test := range tests {
		if got := SubstrAfterLast(test.str, test.find); got != test.want {
			t.Errorf("SubstrAfterLast(%q, %q) = %q; want: %q", test.str, test.find, got, test.want)
		}
	}
}
