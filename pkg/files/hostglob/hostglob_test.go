package hostglob

import (
	"reflect"
	"testing"
)

func TestIsGlob(t *testing.T) {
	testGlob := "myhost[001:003].ci.com"

	if IsGlob(testGlob) != true {
		t.Errorf("failed to return true on %s", testGlob)
	}
}

func TestHostnamePrefix(t *testing.T) {
	testGlob := "myhost[001:003].ci.com"
	want := "myhost"

	if HostnamePrefix(testGlob) != want {
		t.Errorf("hostname prefix did not correctly parse '%s'", want)
	}
}

func TestHostnameSuffix(t *testing.T) {
	testGlob := "myhost[001:003].ci.com"
	want := ".ci.com"

	if HostnameSuffix(testGlob) != want {
		t.Errorf("hostname prefix did not correctly parse '%s'", want)
	}
}

func TestHostRange(t *testing.T) {
	testGlob := "myhost[001:003].ci.com"
	wantStart, wantEnd := "001", "003"

	start, end := HostRange(testGlob)

	if start != wantStart || end != wantEnd {
		t.Errorf("hostranges were not correctly parsed. got start=%s, end=%s", start, end)
	}
}

func TestUncollapse(t *testing.T) {
	tests := map[string]struct {
		input string
		want  []string
	}{
		"with 2 digits": {
			input: "myhost[01:03].ci.com",
			want:  []string{"myhost01.ci.com", "myhost02.ci.com", "myhost03.ci.com"},
		},
		"with 3 digits": {
			input: "myhost[001:003].ci.com",
			want:  []string{"myhost001.ci.com", "myhost002.ci.com", "myhost003.ci.com"},
		},
		"with 4 digits": {
			input: "myhost[0001:0003].ci.com",
			want:  []string{"myhost0001.ci.com", "myhost0002.ci.com", "myhost0003.ci.com"},
		},
	}

	for testName, test := range tests {
		t.Logf("running test %s", testName)
		got, err := Uncollapse(test.input)
		if err != nil {
			t.Error(err)
		}
		if !reflect.DeepEqual(got, test.want) {
			t.Errorf("%s did not uncollapse correctly. got %v, wanted %v", testName, got, test.want)
		}
	}
}
