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
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/client"
	"github.com/masterzen/winrm"
	"github.com/ossan-dev/robocopy/internal/file"
	"github.com/stretchr/testify/require"
)

const windowsImageName = "dockurr/windows"

var (
	winrmClient          *winrm.Client
	dockerClient         *client.Client
	extractorContainerID string
	windowsContainerID   *string
)

func TestRobocopy(t *testing.T) {
	// arrange
	ctx, cancelFunc := context.WithCancel(t.Context())
	t.Cleanup(cancelFunc)
	require.NoError(t, setupFS(), "setupFS")
	// act
	stdOut, stdErr, exitCode, err := winrmClient.RunPSWithContext(ctx, "C:\\Users\\Docker\\Desktop\\Shared\\go-robocopy.exe")
	// assert
	require.NoError(t, err, "RunPSWithContext")
	fmt.Println("stdOut:", stdOut)
	fmt.Println("stdErr:", stdErr)
	require.Zero(t, exitCode, "RunPSWithContext")
	_, err = os.Stat("testdata/target/file.txt")
	require.NoError(t, err, "Stat")
}

func copyBinaryToWindowsContainer(ctx context.Context) error {
	// creation of the extractor container
	containerRes, err := dockerClient.ContainerCreate(ctx, &container.Config{
		Image: "test-my-robocopy",
		Cmd:   []string{"echo", "hello"},
	}, nil, nil, nil, "extractor")
	if err != nil {
		return err
	}
	extractorContainerID = containerRes.ID
	// pull out the binary from the extractor container
	reader, _, err := dockerClient.CopyFromContainer(ctx, containerRes.ID, "/go-robocopy.exe")
	if err != nil {
		return err
	}
	defer reader.Close()
	// copy binary to the target Windows Container
	return dockerClient.CopyToContainer(ctx, *windowsContainerID, "./shared", reader, container.CopyToContainerOptions{})
}

func setupDockerClient(ctx context.Context) (err error) {
	dockerClient, err = client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return err
	}
	dockerClient.NegotiateAPIVersion(ctx)
	return nil
}

func buildDockerImageForRobocopyBinary(ctx context.Context) error {
	// create tarball from source code
	buf := new(bytes.Buffer)
	tw := tar.NewWriter(buf)
	defer tw.Close()
	if err := createTarball(".", tw); err != nil {
		return err
	}
	twReader := bytes.NewReader(buf.Bytes())
	imageBuildRes, err := dockerClient.ImageBuild(ctx, twReader, build.ImageBuildOptions{
		Tags:   []string{"test-my-robocopy"},
		Remove: true,
	})
	if err != nil {
		return err
	}
	defer imageBuildRes.Body.Close()
	if _, err := io.Copy(os.Stdout, imageBuildRes.Body); err != nil {
		return err
	}
	return nil
}

func removeContainerByID(ctx context.Context, containerID string) error {
	return dockerClient.ContainerRemove(ctx, containerID, container.RemoveOptions{})
}

func removeDockerImagesByRepoTags(ctx context.Context, repo, tag string) error {
	filters := filters.NewArgs(filters.Arg("label", fmt.Sprintf("repo=%v", repo)), filters.Arg("label", fmt.Sprintf("tag=%v", tag)))
	images, err := dockerClient.ImageList(ctx, image.ListOptions{Filters: filters})
	if err != nil {
		return err
	}
	for _, img := range images {
		if _, err := dockerClient.ImageRemove(ctx, img.ID, image.RemoveOptions{}); err != nil {
			return err
		}
	}
	return nil
}

func setupFS() error {
	if err := file.FileCreation("testdata/source/file.txt", strings.NewReader(`Hello from Windows in Docker.
This is the file it should be copied by using robocopy.`)); err != nil {
		return err
	}
	if _, err := os.Stat("testdata/target"); err != nil {
		if err := os.MkdirAll("testdata/target", 0700); err != nil {
			return err
		}
		return nil
	}
	if _, err := os.Stat("testdata/target/file.txt"); err == nil {
		if err := os.Remove("testdata/target/file.txt"); err != nil {
			return err
		}
	}
	return nil
}

func getWindowsContainerID(ctx context.Context, dockerClient *client.Client) (*string, error) {
	containers, err := dockerClient.ContainerList(ctx, container.ListOptions{})
	if err != nil {
		return nil, err
	}
	for _, c := range containers {
		if c.Image == windowsImageName {
			return &c.ID, nil
		}
	}
	return nil, fmt.Errorf("no running Windows container found")
}

func getWinrmClient() (*winrm.Client, error) {
	endpoint := winrm.NewEndpoint(
		"localhost",
		5985,
		false,
		false,
		nil,
		nil,
		nil,
		0,
	)
	return winrm.NewClient(endpoint, "Docker", "admin")
}

func createTarball(src string, tw *tar.Writer) error {
	if _, err := os.Stat(src); err != nil {
		return err
	}
	return filepath.WalkDir(src, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if strings.HasPrefix(path, "windows") {
			return nil
		}
		if d.IsDir() {
			return nil
		}
		fileinfo, err := d.Info()
		if err != nil {
			return err
		}
		header, err := tar.FileInfoHeader(fileinfo, path)
		if err != nil {
			return err
		}
		header.Name = path
		if err := tw.WriteHeader(header); err != nil {
			return err
		}
		file, err := os.Open(path)
		if err != nil {
			return err
		}
		defer file.Close()
		if _, err := io.Copy(tw, file); err != nil {
			return err
		}
		return nil
	})
}
