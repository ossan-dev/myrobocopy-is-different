// Package robocopy provides functions to the CLI command
package robocopy

import (
	"fmt"
	"strconv"
	"strings"
)

const (
	lowerBoundRobocopySuccessCode = 0
	upperBoundRobocopySuccessCode = 8
)

// IsError discerns between actual errors and unhappy status codes provided by "robocopy".
// 0	No files were copied. No failure was met. No files were mismatched. The files already exist in the destination directory; so the copy operation was skipped.
// 1	All files were copied successfully.
// 2	There are some additional files in the destination directory that aren't present in the source directory. No files were copied.
// 3	Some files were copied. Additional files were present. No failure was met.
// 5	Some files were copied. Some files were mismatched. No failure was met.
// 6	Additional files and mismatched files exist. No files were copied and no failures were met. Which means that the files already exist in the destination directory.
// 7	Files were copied, a file mismatch was present, and additional files were present.
// 8	Several files didn't copy.
func IsError(potentialError error) error {
	segments := strings.Split(potentialError.Error(), " ")
	resCodeSegment := segments[len(segments)-1]
	resCode, err := strconv.Atoi(resCodeSegment)
	if err != nil {
		return fmt.Errorf("unexpected string value in the last part of the result: %w", err)
	}
	if resCode < lowerBoundRobocopySuccessCode || resCode >= upperBoundRobocopySuccessCode {
		return fmt.Errorf("something went wrong with 'robocopy': %w", potentialError)
	}
	return nil
}
