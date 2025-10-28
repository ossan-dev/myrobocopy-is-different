//go:build integration

package tests

import (
	"context"
	"fmt"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	// setup
	var err error
	ctx, cancelFunc := context.WithCancel(context.Background())
	defer cancelFunc()
	if err := setupDockerClient(ctx); err != nil {
		fmt.Fprintf(os.Stderr, "failed setup Docker client: %v", err.Error())
		os.Exit(1)
	}
	windowsContainerID, err = getWindowsContainerID(ctx, dockerClient)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed checking the Windows container existence: %v", err.Error())
		os.Exit(1)
	}
	if err := buildDockerImageForRobocopyBinary(ctx); err != nil {
		fmt.Fprintf(os.Stderr, "failed building Docker Image: %v", err.Error())
		os.Exit(1)
	}
	if err := copyBinaryToWindowsContainer(ctx); err != nil {
		fmt.Fprintf(os.Stderr, "failed to upload the binary to the Win container: %v", err.Error())
		os.Exit(1)
	}
	winrmClient, err = getWinrmClient()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed getting Winrm client: %v", err.Error())
		os.Exit(1)
	}
	exitCode := m.Run()

	// teardown
	if err := removeDockerImagesByRepoTags(ctx, "test-my-robocopy", "latest"); err != nil {
		fmt.Fprintf(os.Stderr, "fail removing Docker img for 'extractor' container: %v", err.Error())
		os.Exit(1)
	}
	if err := removeContainerByID(ctx, extractorContainerID); err != nil {
		fmt.Fprintf(os.Stderr, "fail removing 'extractor' container: %v", err.Error())
		os.Exit(1)
	}
	if err := dockerClient.Close(); err != nil {
		fmt.Fprintf(os.Stderr, "fail while closing Docker client: %v", err.Error())
		os.Exit(1)
	}
	// teardown
	os.Exit(exitCode)
}
