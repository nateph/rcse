package files

import (
	"testing"
)

func TestParseAndVerifyFilePath(t *testing.T) {
	nonExistentFile := "/temp/shouldnt_exist_because_temp.yaml"

	if _, err := ParseAndVerifyFilePath(nonExistentFile); err == nil {
		t.Error("Didn't recieve error back from function call using non-existant file.")
	}
}
