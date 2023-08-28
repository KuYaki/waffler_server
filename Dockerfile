# syntax=docker/dockerfile:1

FROM golang:1.20 as builder

COPY .. /go/src/waffler/
WORKDIR /go/src/waffler/
RUN CGO_ENABLED=0 go build -a -installsuffix cgo -o waffler ./cmd/api

FROM alpine:3.17.0 as production

RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /go/src/waffler/ ./
CMD ["./waffler"]