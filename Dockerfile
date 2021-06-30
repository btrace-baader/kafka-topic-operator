# Build the manager binary
FROM golang:1.16 as builder

WORKDIR /app
COPY . /app
# cache deps before building and copying source so that we don't need to re-download as much
# and so that source changes don't invalidate our downloaded layer
RUN go mod download

# Build
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 GO111MODULE=on go build -a -o manager

FROM alpine:3.9
WORKDIR /
RUN apk add --no-cache ca-certificates && update-ca-certificates
COPY --from=builder /app/manager .

ENTRYPOINT ["/manager"]
