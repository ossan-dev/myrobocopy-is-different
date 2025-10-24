//go:build integration

package tests

import (
	"archive/tar"
	"bytes"
	"context"
	_ "embed"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/docker/docker/api/types/build"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/stretchr/testify/require"
)

const WindowsImageName = "dockurr/windows"

func getWindowsContainerID(t *testing.T, ctx context.Context, dockerClient *client.Client) *string {
	t.Helper()
	containers, err := dockerClient.ContainerList(ctx, container.ListOptions{})
	require.NoError(t, err, "ContainerList")
	for _, c := range containers {
		if c.Image == WindowsImageName {
			return &c.ID
		}
	}
	return nil
}

func createTarball(t *testing.T, src string, tw *tar.Writer) error {
	t.Helper()
	_, err := os.Stat(src)
	require.NoError(t, err, "Stat")
	return filepath.WalkDir(src, func(path string, d fs.DirEntry, err error) error {
		require.NoError(t, err, "WalkDir")
		if strings.HasPrefix(path, "windows") {
			return nil
		}
		if d.IsDir() {
			return nil
		}
		fileinfo, err := d.Info()
		require.NoError(t, err, "Info")
		header, err := tar.FileInfoHeader(fileinfo, path)
		require.NoError(t, err, "FileInfoHeader")
		header.Name = path
		err = tw.WriteHeader(header)
		require.NoError(t, err, "WriteHeader")
		file, err := os.Open(path)
		require.NoError(t, err, "Open")
		defer file.Close()
		_, err = io.Copy(tw, file)
		require.NoError(t, err, "Copy")
		return nil
	})
}

func TestRobocopy(t *testing.T) {
	// TODO: check for folders existence
	// check for Windows container existence
	ctx, cancelFunc := context.WithCancel(context.Background())
	t.Cleanup(cancelFunc)
	dockerClient, err := client.NewClientWithOpts(client.FromEnv)
	require.NoError(t, err, "NewClientWithOpts")
	dockerClient.NegotiateAPIVersion(ctx)
	windowsContainerID := getWindowsContainerID(t, ctx, dockerClient)
	require.NotNil(t, windowsContainerID, "Windows container is not running. Testing is not possible. Run `make win_setup` and re-test")
	// create tar
	buf := new(bytes.Buffer)
	tw := tar.NewWriter(buf)
	defer tw.Close()
	require.NoError(t, createTarball(t, ".", tw), "createTarball")
	// build and tag Docker image
	twReader := bytes.NewReader(buf.Bytes())
	imageBuildRes, err := dockerClient.ImageBuild(ctx, twReader, build.ImageBuildOptions{
		Tags:   []string{"test-my-robocopy"},
		Remove: true,
	})
	require.NoError(t, err, "ImageBuild")
	defer imageBuildRes.Body.Close()
	// NICETOHAVE: add docker image remove
	_, err = io.Copy(os.Stdout, imageBuildRes.Body)
	require.NoError(t, err, "Copy")
	// container create
	containerRes, err := dockerClient.ContainerCreate(ctx, &container.Config{
		Image: "test-my-robocopy",
		Cmd:   []string{"echo", "hello"},
	}, nil, nil, nil, "extractor")
	require.NoError(t, err, "ContainerCreate")
	fmt.Println(containerRes.ID)
	// container remove
	t.Cleanup(func() {
		require.NoError(t, dockerClient.ContainerRemove(ctx, containerRes.ID, container.RemoveOptions{}), "ContainerRemove")
	})
	// copy from container
	reader, _, err := dockerClient.CopyFromContainer(ctx, containerRes.ID, "/go-robocopy.exe")
	require.NoError(t, err, "CopyFromContainer")
	defer reader.Close()
	// copy to container
	require.NoError(t, dockerClient.CopyToContainer(ctx, *windowsContainerID, "./shared", reader, container.CopyToContainerOptions{}), "CopyToContainer")
}
