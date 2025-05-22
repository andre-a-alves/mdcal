package cmd

import "fmt"

const (
	major = 0
	minor = 4
	patch = 0
)

func getVersion() string {
	return fmt.Sprintf("%d.%d.%d", major, minor, patch)
}
