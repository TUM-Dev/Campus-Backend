FROM golang:alpine as builder

# Install gcc for cgo
RUN apk add build-base

COPY ./server $GOPATH/server/
WORKDIR $GOPATH/server/

RUN go build -o /backend
RUN chmod +x /backend

FROM scratch
COPY --from=builder /backend /backend
CMD ["/backend"]