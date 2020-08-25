package version

import "fmt"

// VERSION indicates which version of the binary is running.
var VERSION = "v0.0.8"

// GITCOMMIT indicates which git hash the binary was built off of
var GITCOMMIT = ""

// Long returns the long representation of version, with semantic version and short git hash (if present)
func Long() string {
	if len(GITCOMMIT) > 0 {
		return fmt.Sprintf("%s-%s", VERSION, GITCOMMIT)
	}
	return VERSION
}
