FROM golang:1.24-alpine AS build

# Install C toolchain
RUN apk add --no-cache \
      gcc \
      musl-dev

WORKDIR /app

# Copy module files
COPY ./go.mod .
COPY ./go.sum .

# Download dependencies
RUN go mod download

# Copy the source code
COPY . .

# Build
RUN CGO_ENABLED=1 go build -o api .

FROM vegardit/osslsigncode:latest-alpine AS osslsigncode

FROM alpine

WORKDIR /app
RUN apk add --no-cache \
      ca-certificates \
      libcrypto3 \
      libssl3 \
      zlib

COPY --from=build /app/api /app/api
COPY --from=osslsigncode /usr/local/bin/osslsigncode /usr/local/bin/osslsigncode

# Create data & config volume
VOLUME ["/app/.data", "/app/config.json"]

# Run the compiled binary
ENTRYPOINT ["/app/api"]
