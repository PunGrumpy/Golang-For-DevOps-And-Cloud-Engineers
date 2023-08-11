# Build Go Project
FROM golang:1.20-alpine AS builder
WORKDIR /oidc-demo
COPY . .
RUN apk add -u -t build-tools curl git && \
    CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o server cmd/server/*.go && \
    apk del build-tools && \
    rm -rf /var/cache/apk/*

# Runtime Container
FROM alpine:latest
WORKDIR /app
RUN apk --no-cache add ca-certificates bash curl
COPY --from=builder /oidc-demo/server /app/
COPY config.yaml /app/config.yaml

ENTRYPOINT [ "/app/oidc-demo-server" ]