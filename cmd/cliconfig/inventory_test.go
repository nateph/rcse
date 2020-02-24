package cliconfig

import (
	"fmt"
	"testing"

	"github.com/spf13/afero"
)

func TestReadInventoryFile(t *testing.T) {
	var yamlExample = []byte(`
name: job to run
hosts:
- myhost.test.com
- example.com
- host.ci.net
other:
  puppet: yes
  updated: no
test_host: true
`)
	fakeFS := afero.NewMemMapFs()
	fakeFS.MkdirAll("/tmp/inv_files/", 0755)
	afero.WriteFile(fakeFS, "/tmp/inv_files/test_inv.yaml", yamlExample, 0644)
	// testFile, err := fakeFS.Stat("/tmp/test_inv.yaml")
	// // if os.IsNotExist(err) {
	// // 	t.Errorf("file \"%s\" does not exist.\n", "/tmp/test_inv.yaml")
	// // }
	todo := ParseAndVerifyFile("/tmp/test_inv.yaml")
	fmt.Println(todo)
	// fakeFileName := fakeFile.Name()
	// fmt.Println(fakeFileName)
	// parsedHosts := ReadInventoryFile(fakeFileName)
	// if !reflect.DeepEqual(parsedHosts, []string{"myhost.test.com", "example.com", "host.ci.net"}) {
	// 	t.Errorf("Hosts list does not match what was expected. Got back: %s", parsedHosts)
	// }
}
