FROM golang:1.25.3-bookworm AS builder

WORKDIR /app

COPY go.mod .
COPY cmd/robocopy/main.go cmd/robocopy/main.go
COPY internal/robocopy/robocopy.go internal/robocopy/robocopy.go

RUN CGO_ENABLED=0 GOOS=windows go build -o go-robocopy.exe cmd/robocopy/main.go

FROM scratch AS final

COPY --from=builder /app/go-robocopy.exe /go-robocopy.exe