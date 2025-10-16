# win_release releases only for Windows by using 'goreleaser'
win_release_ci:
	GOOS=windows GOARCH=arm64 goreleaser release --clean