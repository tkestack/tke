package version

import (
	"reflect"
	"testing"
)

func TestParse(t *testing.T) {
	type parseType struct {
		version string
		val     []string
	}

	vals := []parseType{
		{"1.0.0", []string{"1", "0", "0"}},
		{"1...0.0", []string{"1", "0", "0"}},
		{"0.1.build1004", []string{"0", "1", "build", "1004"}},
		{"0.1+build1004.1", []string{"0", "1", "build", "1004", "1"}},
		{"0.1-1.0", []string{"0", "1", "1", "0"}},
		// {"1.0.1构建日期2014", []string{"1", "0", "1", "构建日期", "2014"}},
	}

	for _, v := range vals {
		verStr, err := Parse(v.version)
		if err != nil {
			t.Fatalf("Couldn't parse version %v: %v", v.version, err)
		}
		if !reflect.DeepEqual(verStr, v.val) {
			t.Fatalf("Parse result is not correct: %v", v.version)
		}
	}

}

func TestCompare(t *testing.T) {
	const (
		gt = iota
		lt
		eq
	)

	var tests = []struct {
		name     string
		v1       string
		v2       string
		expected int
	}{
		{"comp 1", "0.1.0", "0.1.0", eq},
		{"comp 2", "1...0.0", "1.0.0", eq},
		{"comp 3", "1.0-alpha", "1.0-", lt},
		{"comp 4", "1.0+build1", "1.0build1.1", lt},
		{"comp 5", "1.0.build1.1", "1.0build", gt},
		{"comp 6", "1.5.0", "1.5.0", eq},
		{"comp 7", "1.5.1", "1.5.0", gt},
		{"comp 8", "1.6.0", "1.5.1", gt},
	}

	for _, rt := range tests {
		t.Run(rt.name, func(t *testing.T) {
			res := Compare(rt.v1, rt.v2)
			var op int
			if res > 0 {
				op = gt
			} else if res < 0 {
				op = lt
			} else {
				op = eq
			}

			if op != rt.expected {
				t.Errorf(
					"Failed compare %s and %s:\n\texpected: %d\n\t  actual: %d",
					rt.v1,
					rt.v2,
					rt.expected,
					op,
				)
			}
		})
	}
}
