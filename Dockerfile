FROM golang:1.22-alpine AS build

WORKDIR /src
COPY go.mod ./
RUN go mod download
COPY cmd/ ./cmd/
COPY internal/ ./internal/

ENV CGO_ENABLED=0
RUN go build -ldflags "-s -w" -o /out/iptv-proxy ./cmd/iptv-proxy

FROM alpine:3.20 AS certs
RUN apk add --no-cache ca-certificates

FROM scratch
COPY --from=certs /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=build /out/iptv-proxy /iptv-proxy
EXPOSE 8080
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
  CMD ["/iptv-proxy", "healthcheck"]
ENTRYPOINT ["/iptv-proxy"]
