FROM golang:1.22.2 AS builder

WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -a -installsuffix cgo -o process cmd/main.go
RUN chmod +x /app/process



FROM alpine:3.20

COPY --from=builder /app/process /process
# COPY AmazonRootCA1.pem ./AmazonRootCA1.pem
EXPOSE 8080
CMD ["/process"]


