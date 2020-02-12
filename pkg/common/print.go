package common

import (
	"fmt"
)

// PrintHostOutput is a function that can be used by many other functions to
// pretty print a hosts output/other information
func PrintHostOutput(result CommandResult) {
	fmt.Printf("----- %s ----- \n%s\n", result.Host, result.Stdout)
}
