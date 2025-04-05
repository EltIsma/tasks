# Build stage
FROM golang:1.23-alpine AS builder
WORKDIR /task
COPY . .
RUN go mod download
RUN CGO_ENABLED=0 go build -o /task/app /task/cmd/main.go

# Run stage
FROM alpine
WORKDIR /task
COPY --from=builder /task/app .

EXPOSE 8090
ENTRYPOINT [ "./app", "-config", "/etc/task/config.yaml", "-env", "/etc/task/.env"]