FROM golang:1.19.0-alpine3.15 as builder

# Install git + SSL ca certificates.
# Git is required for fetching the dependencies.
# Ca-certificates is required to call HTTPS endpoints.
RUN apk update && apk add --no-cache git ca-certificates tzdata alpine-sdk bash && update-ca-certificates

# Create appuser
RUN adduser -D -g '' appuser
WORKDIR /app

# Copy source code
COPY . .

# bundle version into binary if specified in build-args, dev otherwise.
ARG version=dev
# Compile statically
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -ldflags "-w -extldflags '-static' -X main.Version=${version}" -o /backend main.go
RUN chmod +x /backend

FROM scratch
COPY --from=builder /backend /backend

# Import from builder - needed for running
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /etc/passwd /etc/passwd

# Use an unprivileged user
USER appuser

# Run the main binary
CMD ["/backend"]
