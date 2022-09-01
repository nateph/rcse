package command

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	color "github.com/fatih/color"
	log "github.com/sirupsen/logrus"
	terminal "golang.org/x/term"
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
		green := color.New(color.FgGreen).SprintFunc()
		fmt.Printf(
			"host: %s\ncommand: %s\n",
			green(r.Host),
			r.CommandRan,
		)
		fmt.Printf(
			"stdout:\n%s%s\n",
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
