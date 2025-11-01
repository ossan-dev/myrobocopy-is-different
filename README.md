# Robocopy for Windows developed by a Linux-guy

Robocopy for Windows developed by a Linux-guy.

> TL;DR: it was a bit challenging but fun 😹

In case you're not familiar with the `robocopy` command, these are the semantics.

## Introduction to the Robocopy command

A sample `robocopy` command is like the below one:

```bash
robocopy C:\Users\Docker\Desktop\Shared\source\ C:\Users\Docker\Desktop\Shared\target\ file.txt
```

The syntax is: `robocopy <source folder> <target folder> <filename>`.

## Infrastructure

### Windows Target Machine Setup

After all, it's a Windows-based command... You must have it installed.

> The first time you do this setup, be sure to have the rest of day free 😄. In the meantime, you can enjoy a ☕ or my GitHub profile (LOL).

To have a fully working Windows environment, please run:

```shell
  make win_setup
```

## Build

To build the project, you can use the command `go_build`. It compiles to tool to target a `Windows` machine.

## Run

To run the CLI application, you should follow these steps:

1. Building the CLI tool via the `make go_build` command
2. Connect to the Windows container exposed on port `8006`
3. Open a **MS-DOS** session (or **PowerShell**)
4. `cd C:\Users\Docker\Desktop\Shared`
5. *To let the program creating the default file to be copied, please set the `DEBUG` environment variable in the shell you're using. On Windows: `set DEBUG=true`*
6. run `.\robocopy.exe`
7. check the results

## Test

Testing is a bit of a hack ☠️. We have End-To-End test at the root directory level, meant to fully cover the solution. I came up with this overcomplicated solution for the sake of the testing purposes. I wanted to have an automated test to add to the GitHub Action.  

This is the testing process in general:

1. Instantiating of a Docker client via the Go SDK for Docker
2. Make sure we have the target Windows machine up and running (it's the same we have used above for the manual run, you know it's a bit expensive to recreate another one for testing 🥶)
3. Building of the Docker Image of our System Under Test. We create a distro-less image with only our executable in it
4. Preparing our target Windows machine to be ready to execute the test binary:
    1. run the `extractor` container based on the just created image
    2. copying out the binary from the `extractor` container
    3. copying to the Windows machine this binary
5. Initialization of the `winrm` client (more on this later)
6. Test running
7. Teardown and resources cleanup

### Windows Remote Management setup

To automatically test the solution, we're using the `winrm` tool. To enable this, on the target Windows Machine we're going to use for testing pursposes, you should run these commands there:

1. `winrm config -quiet`
2. `winrm set winrm/config/client @{TrustedHosts="*"}`
3. `winrm set winrm/config/service/auth @{Basic="true"}`
4. `winrm set winrm/config/service @{AllowUnencrypted="true"}`
5. `reg add HKLM\SOFTWARE\Microsoft\Windows\CurrentVersion\Policies\System /v LocalAccountTokenFilterPolicy /t REG_DWORD /d 1 /f`

### Testing Dependencies

To perform the tests, we need these dependencies:

- Docker Daemon: consumed via the Docker SDK for Go
- Windows Target Machine: invoked via `winrm` commands

### Test Command

To start End-To-End testing, you need to run (at the root level directory):

```bash
go test -tags=integration
```

You can also **debug** the End-To-End test by using the profile `debug-integration-tests` of the `.vscode/launch.json` file.

## Release

### Prerequisites

Make sure to have your valid **GITHUB_TOKEN** in the `.env` file. It will be used to publish your releases.

### Goreleaser

To get a smooth process, we use the [goreleaser](https://github.com/goreleaser/goreleaser) tool.
To run a "dry-run" release, issue: `make win_release_dry_run`.  
  
To run a regular release, issue: `make win_release`.
