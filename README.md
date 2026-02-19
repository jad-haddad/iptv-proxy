# IPTV MTV Lebanon Proxy

## Quickstart

```bash
go run ./cmd/iptv-proxy
```

Jellyfin URLs:

- M3U: `http://<server>:8080/lebanon.m3u`
- EPG: `http://<server>:8080/epg.xml`

Docker (public image):

```bash
docker pull ghcr.io/jad-haddad/iptv-proxy:latest
docker run --rm -p 8080:8080 ghcr.io/jad-haddad/iptv-proxy:latest
```

## Overview

This service:

- Fetches the Lebanon M3U playlist and filters to MTV Lebanon only.
- Normalizes the channel to `tvg-id="mtvlebanon.lb"` for Jellyfin matching.
- Fetches a public XMLTV EPG and filters it to MTV Lebanon only.
- Uses ETags for fast 304 responses when nothing changed.

EPG source (explicit): `https://mdag9904.github.io/lebanon-epg/epg.xml`

## Endpoints

- `GET /lebanon.m3u`
- `GET /epg.xml`
- `GET /health` (returns `{"status":"ok"}`)

## Caching and ETags

- Both `/lebanon.m3u` and `/epg.xml` emit `ETag` headers.
- If the client sends `If-None-Match`, the service returns `304 Not Modified` when unchanged.
- EPG fetches are rate-limited in memory by `EPG_REFRESH_SECONDS` (default: 3600 seconds).
- M3U filtering runs only when the upstream M3U changes (via upstream ETag).

## Jellyfin setup

- Add the M3U URL: `http://<server>:8080/lebanon.m3u`
- Add the EPG URL: `http://<server>:8080/epg.xml`

The service ensures `tvg-id="mtvlebanon.lb"` so Jellyfin can match the guide automatically.

## Configuration

Environment variables:

- `M3U_URL` (default: `https://iptv-org.github.io/iptv/countries/lb.m3u`)
- `EPG_URL` (default: `https://mdag9904.github.io/lebanon-epg/epg.xml`)
- `MTV_REGEX` (default: `(?i)\bmtv\b.*\blebanon\b|mtv\s*lebanon|mtvlebanon`)
- `MTV_TVG_ID` (default: `mtvlebanon.lb`)
- `MTV_TVG_NAME` (default: `MTV Lebanon`)
- `EPG_REFRESH_SECONDS` (default: `3600`)
- `REQUEST_TIMEOUT_SECONDS` (default: `15`)

## Docker

Build and run locally:

```bash
docker build -t iptv-proxy .
docker run --rm -p 8080:8080 iptv-proxy
```

Pull and run the public image:

```bash
docker pull ghcr.io/jad-haddad/iptv-proxy:latest
docker run --rm -p 8080:8080 ghcr.io/jad-haddad/iptv-proxy:latest
```

### Docker Compose

```yaml
services:
  iptv-proxy:
    image: ghcr.io/jad-haddad/iptv-proxy:latest
    ports:
      - "8080:8080"
    healthcheck:
      test: ["CMD", "/iptv-proxy", "healthcheck"]
      interval: 30s
      timeout: 3s
      start_period: 5s
      retries: 3
```

## Validation

Run the integration check (validates 200 responses and ETag 304 behavior):

```bash
go run ./scripts/check_endpoints.go http://127.0.0.1:8080
```

## Troubleshooting

- Port already in use: choose another host port in the Docker run mapping.
- No EPG data: verify the upstream EPG URL is reachable.
- MTV not matched: adjust `MTV_REGEX` to match new playlist metadata.
