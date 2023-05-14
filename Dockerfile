FROM golang:1.19-alpine AS Builder
WORKDIR /app

# Download Depedency
COPY go.* ./
RUN go mod tidy

# Copy Source Code to Container
COPY . ./

# Build binary
RUN go build -v -o load-testing-proxy-server

FROM alpine:3.9

# Copy binary from build stage
COPY --from=Builder /app/load-testing-proxy-server /load-testing-proxy-server
COPY --from=Builder /app/serviceAccountKey.json /serviceAccountKey.json

# Run app
ENTRYPOINT ["./load-testing-proxy-server"]
