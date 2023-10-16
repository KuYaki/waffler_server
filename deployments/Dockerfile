# syntax=docker/dockerfile:1

FROM golang:1.21.1 as builder

WORKDIR /go/src/waffler/
COPY go.mod ./
COPY go.sum ./
RUN go mod download
COPY .. .
RUN CGO_ENABLED=0 go build -a -installsuffix cgo -o waffler ./cmd/api

FROM alpine:3.17.0 as production

RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /go/src/waffler/ ./
CMD ["./waffler"]