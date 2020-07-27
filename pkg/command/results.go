package command

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/prometheus/common/log"
	"golang.org/x/crypto/ssh/terminal"
	"gopkg.in/yaml.v2"
)

// Result contains various information about what a command returned.
type Result struct {
	// The command that was ran.
	CommandRan string `json:"command" yaml:"command"`
	// stdout from the command.
	Stdout string `json:"result" yaml:"result"`
	// Host command ran on
	Host string `json:"host" yaml:"host"`
}

// PrintHostOutput formats the host and stdout nicely.
func (r *Result) PrintHostOutput(format string) {
	terminalWidth := getTerminalWidth()
	switch format {
	case "text":
		fmt.Printf(
			"host: %s\ncommand: %s\nstdout:\n%s%s\n",
			r.Host,
			r.CommandRan,
			r.Stdout,
			strings.Repeat("-", terminalWidth),
		)
	case "json":
		jsonData, err := json.Marshal(&r)
		if err != nil {
			return
		}
		fmt.Println(string(jsonData))
	case "yaml":
		yamlData, err := yaml.Marshal(&r)
		if err != nil {
			return
		}
		fmt.Println(string(yamlData))
	}
}

func getTerminalWidth() int {
	termID := int(os.Stdout.Fd())
	width, _, err := terminal.GetSize(termID)
	if err != nil {
		log.Error(err)
	}
	return width
}
