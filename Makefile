include .env

# win_setup runs the docker 
win_setup:
	docker compose up --no-recreate -d

# go_build prepares the test binary
go_build:
	GOOS=windows go build -o testdata/robocopy.exe cmd/robocopy/main.go

# win_build creates the binary for Windows
win_build: export_token
	export GITHUB_TOKEN=$(GITHUB_TOKEN) && GOOS=windows GOARCH=arm64 goreleaser build --snapshot --clean

# win_release_dry_run is the 'no-op' version of the 'win_release' command
win_release_dry_run:
	export GITHUB_TOKEN=$(GITHUB_TOKEN) && GOOS=windows GOARCH=arm64 goreleaser release --snapshot --clean --skip=publish

# win_release releases only for Windows by using 'goreleaser'
win_release: export_token
	export GITHUB_TOKEN=$(GITHUB_TOKEN) && GOOS=windows GOARCH=arm64 goreleaser release --clean
