package cmd

import "fmt"

const (
	major = 0
	minor = 2
	patch = 2
)

func getVersion() string {
	return fmt.Sprintf("%d.%d.%d", major, minor, patch)
}
