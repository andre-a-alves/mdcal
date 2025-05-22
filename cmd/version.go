package cmd

import "fmt"

const (
	major = 0
	minor = 1
	patch = 1
)

func getVersion() string {
	return fmt.Sprintf("%d.%d.%d", major, minor, patch)
}
