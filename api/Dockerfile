# Build arguments
ARG API_ALPINE_VERSION=3.19.1
ARG API_ARCH=amd64

# STEP 1 build executable binary
FROM golang:alpine AS builder
RUN adduser -D -g '' appuser && \
    mkdir /app
WORKDIR /app
# copy the source code
COPY . .
# build the binary
RUN go mod download && \
    CGO_ENABLED=0 GOOS=linux GOARCH=${API_ARCH} go build -a -installsuffix cgo -ldflags="-w -s" -o /go/bin/app . && \
    chmod 0555 /go/bin/app

# STEP 2 build a small image
FROM alpine:${API_ALPINE_VERSION}
WORKDIR /app
COPY --from=builder /etc/passwd /etc/passwd
USER appuser
# Copy our static executable
COPY --from=builder /go/bin/app ./app_binary
CMD ["./app_binary"]