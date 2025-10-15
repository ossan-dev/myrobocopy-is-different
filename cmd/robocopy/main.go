package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"runtime"

	"github.com/ossan-dev/robocopy/internal/robocopy"
)

var (
	sourceDir string
	targetDir string
	filename  string
)

func init() {
	flag.StringVar(&sourceDir, "sourceDir", "C:\\Users\\Docker\\Desktop\\Shared\\source\\", "source directory where the file lives")
	flag.StringVar(&targetDir, "targetDir", "C:\\Users\\Docker\\Desktop\\Shared\\target\\", "target directory where the file will be copied")
	flag.StringVar(&filename, "filename", "file.txt", "name of the file to copy")
}

func main() {
	if runtime.GOOS != "windows" {
		fmt.Fprintln(os.Stdout, "not running on Windows. Running on:", runtime.GOOS)
		return
	}
	flag.Parse()
	cmdName := "cmd.exe"
	cmdArgs := []string{cmdName, "/c", "robocopy", sourceDir, targetDir, filename} // "/c carries out the command and then terminates"
	cmd := exec.Command(cmdName, cmdArgs...)
	result, err := cmd.Output()
	if err != nil && robocopy.AssessError(err) != nil {
		fmt.Fprintf(os.Stderr, "robocopy err: %v", robocopy.AssessError(err).Error())
	}
	fmt.Fprintln(os.Stdout, string(result))
}
