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

// func TestLoadReader(t *testing.T) {
// 	var TestFs = afero.NewMemMapFs()
// 	var yamlExample = []byte(`
// hosts:
// - myhost.test.com
// - example.com
// - host.ci.net
// test_host: true
// `)
// 	// f, err := afero.TempFile(TestFs, "", "test_inv")
// 	TestFs.MkdirAll("test_folder", 0755)
// 	afero.WriteFile(TestFs, "test_folder/test_inv", yamlExample, 0644)
// 	myInv, err := LoadReader("test_folder/test_inv")
// 	if err != nil {
// 		return err
// 	}
// 	fmt.Println(myInv)
// }
