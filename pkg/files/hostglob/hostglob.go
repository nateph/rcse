package hostglob

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// IsGlob takes a line, likely from a file, and returns whether or not
// it is a range of hosts
func IsGlob(line string) bool {
	r := regexp.MustCompile(`\[\d+:\d+\]`)
	return r.MatchString(line)
}

// Uncollapse will return a slice of all hosts from a globbed host range
func Uncollapse(glob string) ([]string, error) {
	prefix := HostnamePrefix(glob)
	suffix := HostnameSuffix(glob)

	startRange, endRange := HostRange(glob)

	// We can use start or end value here, they should be the same
	// This is used for formatting the hostname
	digits := len(startRange)

	start, err := strconv.Atoi(startRange)
	if err != nil {
		return []string{}, err
	}

	end, err := strconv.Atoi(endRange)
	if err != nil {
		return []string{}, err
	}

	var hosts []string

	// todo: find a way to clean this up. right now I can't figure out how to
	// insert a variable into a decimal counted formatted string
	switch digits {
	case 2:
		for i := start; i <= end; i++ {
			hosts = append(hosts, fmt.Sprintf("%s%02d%s", prefix, i, suffix))
		}
	case 3:
		for i := start; i <= end; i++ {
			hosts = append(hosts, fmt.Sprintf("%s%03d%s", prefix, i, suffix))
		}
	case 4:
		for i := start; i <= end; i++ {
			hosts = append(hosts, fmt.Sprintf("%s%04d%s", prefix, i, suffix))
		}
	default:
		for i := start; i <= end; i++ {
			hosts = append(hosts, fmt.Sprintf("%s%d%s", prefix, i, suffix))
		}
	}

	return hosts, nil
}

// HostnamePrefix will return the named prefix of the host, before the [:] range specifier
// for example, passing in myhost[001:008].ci.com would return "myhost"
func HostnamePrefix(glob string) string {
	bracket := strings.Index(glob, "[")
	prefix := glob[:bracket]
	return prefix
}

// HostnameSuffix will return the named suffix of the host, after the [:] range specifier
// for example, passing in myhost[001:008].ci.com would return ".ci.com"
func HostnameSuffix(glob string) string {
	bracket := strings.Index(glob, "]")
	suffix := glob[bracket+1:]
	return suffix
}

// HostRange finds where the range starts and ends
func HostRange(glob string) (string, string) {
	colon := strings.Index(glob, ":")
	startingBracket := strings.Index(glob, "[")
	endingBracket := strings.Index(glob, "]")

	start := glob[startingBracket+1 : colon]

	end := glob[colon+1 : endingBracket]

	return start, end
}
