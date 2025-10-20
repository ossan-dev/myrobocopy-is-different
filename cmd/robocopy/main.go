package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	"github.com/ossan-dev/robocopy/internal/robocopy"
)

var (
	sourceDir string
	targetDir string
	filename  string
)

const (
	sourceDirDebug = "C:\\Users\\Docker\\Desktop\\Shared\\source\\"
	targetDirDebug = "C:\\Users\\Docker\\Desktop\\Shared\\target\\"
	filenameDebug  = "file.txt"
)

func init() {
	flag.StringVar(&sourceDir, "sourceDir", sourceDirDebug, "source directory where the file lives")
	flag.StringVar(&targetDir, "targetDir", targetDirDebug, "target directory where the file will be copied")
	flag.StringVar(&filename, "filename", filenameDebug, "name of the file to copy")
}

func main() {
	if runtime.GOOS != "windows" {
		fmt.Fprintln(os.Stdout, "not running on Windows. Running on:", runtime.GOOS)
		return
	}

	// debug section
	if os.Getenv("DEBUG") == "true" {
		file, err := os.Create(filepath.Join(sourceDirDebug, filenameDebug))
		if err != nil {
			fmt.Fprintf(os.Stderr, "debug mode: failed file creation: %v", err.Error())
			return
		}
		defer file.Close()
		if _, err := file.WriteString(`Hello from Windows in Docker.
This is the file it should be copied by using robocopy.`); err != nil {
			fmt.Fprintf(os.Stderr, "debug mode: failed file content writing: %v", err.Error())
			return
		}
	}

	flag.Parse()
	cmd := &exec.Cmd{
		Path:   "C:\\Windows\\System32\\Robocopy.exe",
		Args:   []string{"robocopy", sourceDir, targetDir, filename}, // "/c carries out the command and then terminates"
		Stdout: os.Stdout,
		Stderr: os.Stdout,
	}
	if err := cmd.Start(); err != nil {
		fmt.Fprintf(os.Stderr, "failed starting robocopy: %v", err.Error())
		return
	}
	if err := cmd.Wait(); err != nil {
		if robocopy.AssessError(err) != nil {
			fmt.Fprintf(os.Stderr, "robocopy err: %v", robocopy.AssessError(err).Error())
		}
	}
}
