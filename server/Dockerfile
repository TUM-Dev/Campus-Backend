FROM golang:1.22-alpine3.19 as builder

# Install git + SSL ca certificates.
# Git is required for fetching the dependencies.
# Ca-certificates is required to call HTTPS endpoints.
RUN apk update && apk add --no-cache git ca-certificates tzdata alpine-sdk bash && update-ca-certificates

# Create appuser
RUN adduser -D -g '' appuser
WORKDIR /app

# Copy go mod and sum files and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY api/tumdev/campus_backend.swagger.json ./swagger/swagger.json
COPY . .

# bundle version into binary if specified in build-args, dev otherwise.
ARG version=dev
# Compile statically
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -ldflags "-w -extldflags '-static' -X main.Version=${version}" -o /backend main.go
RUN chmod +x /backend

FROM alpine:latest
COPY --from=builder /backend /backend

# Import from builder - needed for running
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /etc/passwd /etc/passwd

# make sure that /Storage is mounted with the correct access permissions even in docker-compose
RUN mkdir /Storage
RUN chown appuser /Storage
VOLUME /Storage

# Use an unprivileged user
USER appuser

# metadata
EXPOSE 50051
HEALTHCHECK  --interval=5m --timeout=3s \
  CMD wget --no-verbose --tries=1 --spider http://localhost:50051/health && exit 0 || exit 1

# Run the main binary
CMD ["/backend"]
