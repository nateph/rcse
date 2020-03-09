package command

import (
	"encoding/json"
	"fmt"

	"gopkg.in/yaml.v2"
)

// Results contains all results for a host list
type Results struct {
	Result []*Result `json:"results"`
}

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
	switch format {
	case "text":
		fmt.Printf("----- %s -----\n%s\n\n%s\n", r.Host, r.CommandRan, r.Stdout)
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

// func handleJSON() {

// }
