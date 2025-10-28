FROM golang:1.25.3-bookworm AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY cmd/robocopy/main.go cmd/robocopy/main.go
COPY internal/robocopy/robocopy.go internal/robocopy/robocopy.go
COPY internal/file/local_fs_setup.go internal/file/local_fs_setup.go

RUN CGO_ENABLED=0 GOOS=windows go build -o go-robocopy.exe cmd/robocopy/main.go

FROM scratch AS final

COPY --from=builder /app/go-robocopy.exe /go-robocopy.exe