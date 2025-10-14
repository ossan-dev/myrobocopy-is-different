package main

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
)

func main() {
	if runtime.GOOS != "windows" {
		fmt.Println("not running on Windows. Running on:", runtime.GOOS)
		return
	}
	cmdName := "cmd.exe"
	cmdArgs := []string{"/c", "dir"} // "/c carries out the command and then terminates"
	cmd := exec.Command(cmdName, cmdArgs...)
	var dirs []byte
	var err error
	if dirs, err = cmd.Output(); err != nil {
		fmt.Fprintln(os.Stderr, "error while executing 'dir' cmd:", err.Error())
		return
	}
	fmt.Println("dir:")
	fmt.Println(string(dirs))
}
