# BUILDER
FROM golang:1.13 AS builder
WORKDIR /app
COPY go.mod go.sum .
RUN go mod download
COPY main.go .
RUN go install -v

# FINAL IMAGE
FROM gcr.io/distroless/base
COPY --from=builder /go/bin/sha3 /bin/sha3
ENTRYPOINT ["/bin/sha3"]
