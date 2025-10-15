// Package robocopy provides functions to the CLI command
package robocopy

import (
	"fmt"
	"os/exec"
)

const (
	fileCopiedStatus = 1
)

// AssessError discerns between actual errors and unhappy status codes provided by "robocopy".
//
// Robocopy will return "1" as exit status code when it copies the file. This is interpreted as an error by the OS.
// When we get an error, we have to discern between a "false" error (the file has been copied) and an actual error happened for whatever reason.
func AssessError(cmdResult error) error {
	exitErr, ok := cmdResult.(*exec.ExitError)
	if !ok {
		return fmt.Errorf("generic OS err: %w", cmdResult)
	}
	if exitErr.ExitCode() == fileCopiedStatus {
		return nil
	}
	return exitErr
}
