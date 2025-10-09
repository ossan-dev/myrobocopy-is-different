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
	cmd := exec.Command("dir")
	dirs, err := cmd.Output()
	if err != nil {
		fmt.Fprintln(os.Stderr, "error while executing 'dir' cmd:", err.Error())
		return
	}
	fmt.Println("dir:")
	fmt.Println(dirs)
}
