# ---- Builder Stage ----
FROM golang:1.25-alpine AS builder

WORKDIR /build

# Cache dependency layer separately
COPY go.mod go.sum ./
ENV GOPROXY=https://mirrors.tencent.com/go/
RUN go mod download

# Copy source and build a statically-linked binary
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o gin-demo ./cmd/main.go

# ---- Runtime Stage ----
FROM alpine:3.21

# OIDC requires HTTPS; install CA certificates
RUN apk add --no-cache ca-certificates

WORKDIR /app
COPY --from=builder /build/gin-demo .

EXPOSE 8080

ENTRYPOINT ["./gin-demo"]
