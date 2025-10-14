package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"runtime"
)

const falsePositiveErr = "exit status 1"

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
	if err := cmd.Run(); err != nil && err.Error() != falsePositiveErr {
		fmt.Fprintln(os.Stdout, err.Error(), ":", cmdErr.String())
		return
	}
	fmt.Fprintln(os.Stdout, "Result:", cmdOutput.String())
}
