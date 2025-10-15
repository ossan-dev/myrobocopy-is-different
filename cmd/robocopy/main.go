package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"runtime"

	"github.com/ossan-dev/robocopy/internal/robocopy"
)

func main() {
	if runtime.GOOS != "windows" {
		fmt.Fprintln(os.Stdout, "not running on Windows. Running on:", runtime.GOOS)
		return
	}
	var cmdErr bytes.Buffer
	var cmdOutput bytes.Buffer
	cmdName := "cmd.exe"
	cmdArgs := []string{cmdName, "/c", "robocopy", "C:\\Users\\Docker\\Desktop\\Shared\\source\\", "C:\\Users\\Docker\\Desktop\\Shared\\target\\", "file.txt"} // "/c carries out the command and then terminates"
	cmd := exec.Command(cmdName, cmdArgs...)
	cmd.Stderr = &cmdErr
	cmd.Stdout = &cmdOutput
	if potentialErr := cmd.Run(); potentialErr != nil {
		if err := robocopy.IsError(potentialErr); err != nil {
			fmt.Fprintf(os.Stdout, "cli issue: %v", err.Error())
			return
		}
	}
	fmt.Fprintln(os.Stdout, cmdOutput.String())
}
