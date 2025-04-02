# 🔧 Build bosqichi (stage 1)
FROM golang:1.23 AS builder

WORKDIR /app

COPY ../go.mod ./
COPY ../go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o app main.go

FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

COPY --from=builder /app/app .

COPY storage/keys/ ./storage/keys/


CMD ["./app"]
