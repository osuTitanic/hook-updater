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

FROM alpine

WORKDIR /app
COPY --from=build /app/api /app/api

# Create data & config volume
VOLUME ["/app/.data", "/app/config.json"]

# Run the compiled binary
ENTRYPOINT ["/app/api"]