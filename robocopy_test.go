//go:build integration

package tests

import (
	"archive/tar"
	"bytes"
	"context"
	_ "embed"
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

func isWindowsContainerRunning(t *testing.T, ctx context.Context, dockerClient *client.Client) bool {
	t.Helper()
	containers, err := dockerClient.ContainerList(ctx, container.ListOptions{})
	require.NoError(t, err, "ContainerList")
	for _, c := range containers {
		if c.Image == WindowsImageName {
			return true
		}
	}
	return false
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
	// check for Windows container existence
	ctx, cancelFunc := context.WithCancel(context.Background())
	t.Cleanup(cancelFunc)
	dockerClient, err := client.NewClientWithOpts(client.FromEnv)
	require.NoError(t, err, "NewClientWithOpts")
	dockerClient.NegotiateAPIVersion(ctx)
	require.True(t, isWindowsContainerRunning(t, ctx, dockerClient), "Windows container is not running. Testing is not possible. Run `make win_setup` and re-test")
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
	_, err = io.Copy(os.Stdout, imageBuildRes.Body)
	require.NoError(t, err, "Copy")
}
