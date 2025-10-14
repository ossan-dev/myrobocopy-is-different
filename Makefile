include .env

# win_setup runs the docker 
win_setup:
	docker compose up --no-recreate -d

# go_build prepares the test binary
go_build:
	GOOS=windows go build -o testdata/robocopy.exe cmd/robocopy/main.go

# export_token sets the GITHUB_TOKEN in the shell
export_token:
	export GITHUB_TOKEN=$(GITHUB_TOKEN)

# win_build creates the binary for Windows
win_build: export_token
	GOOS=windows GOARCH=arm64 goreleaser build --snapshot --clean

# win_release releases only for Windows by using 'goreleaser'
win_release: export_token
	GOOS=windows GOARCH=arm64 goreleaser release --snapshot --clean --skip=publish	