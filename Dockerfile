FROM golang:1.22.11-bookworm AS builder
WORKDIR /go/src/builder
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o build/grpc main.go

# https://github.com/GoogleContainerTools/distroless
FROM alpine:3.19

# Copy the server binary
COPY --from=builder /go/src/builder/build/grpc /app/server
COPY docker-entrypoint.sh /app/entrypoint.sh
WORKDIR /app
Expose 8080
CMD ["/bin/sh", "./entrypoint.sh"]