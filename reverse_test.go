package main

import (
	"testing"
	"unicode/utf8"
)

func TestReverse(t *testing.T) {
	testcases := []struct {
		in, want string
	}{
		{in: "The quick brown fox jumps over the lazy dog", want: "god yzal eht revo spmuj xof nworb kciuq ehT"},
		{in: " ", want: " "},
		{in: "!12345", want: "54321!"},
		{in: "和额", want: "额和"},
	}

	for _, tc := range testcases {
		res, err := Reverse(tc.in)
		if err != nil {
			return
		}
		if res != tc.want {
			t.Errorf("Reverse(%q) = %q; want %q", tc.in, res, tc.want)
		}
	}
}

func FuzzReverse(f *testing.F) {
	testcases := []string{"The quick brown fox jumps over the lazy dog", " ", "!12345"}
	for _, testcase := range testcases {
		f.Add(testcase)
	}
	f.Fuzz(func(t *testing.T, s string) {
		res, err1 := Reverse(s)
		if err1 != nil {
			return
		}
		double, err2 := Reverse(res)
		if err2 != nil {
			return
		}
		if s != double {
			t.Errorf("Reverse(%q) = %q; want %q", s, res, double)
		}
		if utf8.ValidString(s) && !utf8.ValidString(res) {
			t.Errorf("Valid utf8 string: %q", res)
		}
	})
}
