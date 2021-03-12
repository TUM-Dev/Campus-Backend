FROM golang:alpine as builder

COPY ./server $GOPATH/server/
WORKDIR $GOPATH/server/

RUN go build -o /backend
RUN chmod +x /backend

FROM alpine:latest
COPY --from=builder /backend /backend

CMD ["/go-direct"]