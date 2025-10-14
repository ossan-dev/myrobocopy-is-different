# Robocopy for Windows developed by a Linux-guy

Robocopy for Windows developed by a Linux-guy.  

In case you're not familiar with the `robocopy` command, these are the semantics.

## Command Syntax

A sample `robocopy` command is like the below one:

```bash
robocopy C:\Users\Docker\Desktop\Shared\source\ C:\Users\Docker\Desktop\Shared\target\ file.txt
```

The syntax is: `robocopy <source folder> <target folder> <filename>`.

## Windows Setup

After all, it's a Windows command... You must have it installed.

> The first you download the Windows Docker image, be sure to have the rest of day free 😄. In the meantime, you can enjoy a ☕ or my GitHub profile (LOL).

To have a fully working Windows environment, please run:

```shell
  make win_setup
```

## Build

To build the project, you can use the command `go_build`. It compiles to tool to target a `windows` machine.

## Test

At the moment, the tests are done manually. The steps are:

1. Building the CLI tool via the `make go_build` command
2. Connect to the Windows container exposed on port `8006`
3. Open a **Windows PowerShell**
4. `cd C:\Users\Docker\Desktop\Shared`
5. run `.\robocopy.exe`
6. check the results

## Release

### Prerequisites

Make sure to have your valid **GITHUB_TOKEN** in the `.env` file. It will be used to publish your releases.

### Goreleaser

To get a smooth process, we use the [goreleaser](https://github.com/goreleaser/goreleaser) tool.
To run a "dry-run" release, issue: `make win_release_dry_run`.  
  
To run a regular release, issue: `make win_release`.
