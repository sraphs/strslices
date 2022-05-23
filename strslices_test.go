package strslices

import (
	"fmt"
	"reflect"
	"regexp"
	"strings"
	"testing"
)

func Example_regexp() {
	a := []string{
		"10.0.0.0", "FOO", "1000::1", "fd80::4%eth0", "BAR", "192.168.1.300",
		"172.12.0.1", "fc00::5000",
	}
	re := regexp.MustCompile("^[A-Z]+$")
	fmt.Println(Filter(nil, a, re.MatchString))
	// Output: [FOO BAR]
}

func ExampleFilter_empty() {
	s := []string{"FOO", "", "", "BAR", "fd80::8888", ""}
	fmt.Println(Filter(nil, s, func(s string) bool { return s != "" }))
	// Output: [FOO BAR fd80::8888]
}

func TestFilter(t *testing.T) {
	testCases := []struct {
		validator func(string) bool
		input     []string
		res       []string
	}{
		// Filter all
		{
			regexp.MustCompile("zzz.*").MatchString,
			[]string{"FOO", "1000::", "BAR", "fd80::8888"},
			nil,
		},
		// Filter none
		{
			regexp.MustCompile(".*").MatchString,
			[]string{"FOO", "1000::", "BAR", "fd80::8888"},
			[]string{"FOO", "1000::", "BAR", "fd80::8888"},
		},
		// Filter some
		{
			func(s string) bool { return strings.Contains(s, ".") },
			[]string{"10.0.0.0", "1000::", "8.8.8.8", "fd80::8888"},
			[]string{"10.0.0.0", "8.8.8.8"},
		},
	}
	for i, tc := range testCases {
		res := Filter(nil, tc.input, tc.validator)
		if !reflect.DeepEqual(tc.res, res) {
			t.Errorf("TC %d: %v expected %v", i, res, tc.res)
		}
	}
}

func TestFilterInplace(t *testing.T) {
	testCases := []struct {
		validator  func(string) bool
		input      []string
		res        []string
		inputAfter []string
	}{
		// Filter all
		{
			func(s string) bool { return s != "" },
			[]string{"FOO", "", "", "BAR", "fd80::8888", ""},
			[]string{"FOO", "BAR", "fd80::8888"},
			[]string{"FOO", "BAR", "fd80::8888", "BAR", "fd80::8888", ""},
		},
	}
	for i, tc := range testCases {
		res := Filter(tc.input[:0], tc.input, tc.validator)
		if !reflect.DeepEqual(tc.res, res) {
			t.Errorf("TC %d: %v expected %v", i, res, tc.res)
		}
		if !reflect.DeepEqual(tc.input, tc.inputAfter) {
			t.Errorf("TC %d: mutated input %v expected %v", i, tc.input, tc.inputAfter)
		}
	}
}

func TestContains(t *testing.T) {
	testCases := []struct {
		input []string
		what  string
		res   bool
	}{
		// Contains special case
		{
			nil,
			"",
			false,
		},
		// Contains
		{
			[]string{"FOO", "BAR", ""},
			"",
			true,
		},
		// Not Contains
		{
			[]string{"FOO", "BAR", ""},
			"NOPE",
			false,
		},
	}
	for i, tc := range testCases {
		res := Contains(tc.input, tc.what)
		if res != tc.res {
			t.Errorf("TC %d: %v expected %v", i, res, tc.res)
		}
	}
}

func TestClone(t *testing.T) {
	testCases := []struct {
		input []string
		res   []string
	}{
		{
			nil,
			nil,
		},
		{
			[]string{},
			[]string{},
		},
		{
			[]string{"", "FOO", "BAR", ""},
			[]string{"", "FOO", "BAR", ""},
		},
	}
	for i, tc := range testCases {
		res := Clone(tc.input)
		if !reflect.DeepEqual(tc.res, res) {
			t.Errorf("TC %d: %v expected %v", i, res, tc.res)
		}
		if len(tc.input) > 0 {
			tc.input[0] = "NOPE"
			if tc.res[0] == "NOPE" {
				t.Errorf("TC %d: Clone is not cloned", i)
			}
		}
	}
}
