FROM golang:1 AS build-env

# Install ca-certificates for Go module downloads
RUN apt-get update && apt-get install -y ca-certificates && rm -rf /var/lib/apt/lists/*

WORKDIR /go/src/github.com/sstreichan/ztdns
# Copy go.mod and go.sum first for better caching
COPY go.mod go.sum ./

# Install dependencies
RUN go mod download

# Add source
COPY . .

# Build static binary
RUN CGO_ENABLED=0 GOOS=linux go install -v ./...

FROM alpine

# We need to add ca-certificates in order to make HTTPS API calls
RUN apk update && apk add ca-certificates && rm -rf /var/cache/apk/*

WORKDIR /app
# Copy binary
COPY --from=build-env /go/bin/ztdns .

ENTRYPOINT ["./ztdns", "--debug"]
CMD ["server"]
EXPOSE 53/udp
