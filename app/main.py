import hashlib
import time

import httpx
from fastapi import FastAPI, Request, Response

from app.cache import EPGCache, M3UCache
from app.config import (
    EPG_REFRESH_SECONDS,
    EPG_URL,
    MTV_REGEX,
    MTV_TVG_ID,
    MTV_TVG_NAME,
    M3U_URL,
    REQUEST_TIMEOUT_SECONDS,
)
from app.epg import filter_epg
from app.m3u import filter_m3u


app = FastAPI(title="IPTV MTV Lebanon Proxy")
client: httpx.AsyncClient | None = None
cache_m3u = M3UCache()
cache_epg = EPGCache()


@app.on_event("startup")
async def _startup() -> None:
    global client
    client = httpx.AsyncClient(timeout=REQUEST_TIMEOUT_SECONDS)


@app.on_event("shutdown")
async def _shutdown() -> None:
    if client:
        await client.aclose()


def _etag_for(content: bytes) -> str:
    digest = hashlib.sha256(content).hexdigest()
    return f'"{digest}"'


def _cache_headers(etag: str) -> dict[str, str]:
    return {"ETag": etag, "Cache-Control": "no-cache"}


def _headers_from_cache(etag: str | None, body: bytes | None) -> dict[str, str]:
    if etag:
        return _cache_headers(etag)
    if body is None:
        return {"Cache-Control": "no-cache"}
    return _cache_headers(_etag_for(body))


@app.get("/health")
async def health() -> dict[str, str]:
    return {"status": "ok"}


@app.get("/lebanon.m3u")
async def get_m3u(request: Request) -> Response:
    assert client is not None
    if (
        cache_m3u.filtered_body
        and request.headers.get("if-none-match") == cache_m3u.filtered_etag
    ):
        return Response(status_code=304)

    async with cache_m3u.lock:
        if (
            cache_m3u.filtered_body
            and request.headers.get("if-none-match") == cache_m3u.filtered_etag
        ):
            return Response(status_code=304)

        headers: dict[str, str] = {}
        if cache_m3u.upstream_etag:
            headers["If-None-Match"] = cache_m3u.upstream_etag

        resp = await client.get(M3U_URL, headers=headers)
        if resp.status_code == 304:
            if cache_m3u.filtered_body:
                return Response(
                    content=cache_m3u.filtered_body,
                    media_type="application/vnd.apple.mpegurl",
                    headers=_headers_from_cache(
                        cache_m3u.filtered_etag, cache_m3u.filtered_body
                    ),
                )
            resp = await client.get(M3U_URL)
        if resp.status_code != 200:
            return Response(status_code=502, content="Upstream M3U unavailable")

        cache_m3u.upstream_etag = resp.headers.get("ETag")
        filtered = filter_m3u(resp.text, MTV_REGEX, MTV_TVG_ID, MTV_TVG_NAME)
        body = filtered.encode("utf-8")
        cache_m3u.filtered_body = body
        cache_m3u.filtered_etag = _etag_for(body)

        if request.headers.get("if-none-match") == cache_m3u.filtered_etag:
            return Response(status_code=304)

        return Response(
            content=body,
            media_type="application/vnd.apple.mpegurl",
            headers=_headers_from_cache(cache_m3u.filtered_etag, body),
        )


@app.get("/epg.xml")
async def get_epg(request: Request) -> Response:
    assert client is not None
    now = time.time()
    if (
        cache_epg.filtered_body
        and now - cache_epg.last_fetch < EPG_REFRESH_SECONDS
        and request.headers.get("if-none-match") == cache_epg.filtered_etag
    ):
        return Response(status_code=304)

    async with cache_epg.lock:
        now = time.time()
        if cache_epg.filtered_body and now - cache_epg.last_fetch < EPG_REFRESH_SECONDS:
            if request.headers.get("if-none-match") == cache_epg.filtered_etag:
                return Response(status_code=304)
            return Response(
                content=cache_epg.filtered_body,
                media_type="application/xml",
                headers=_headers_from_cache(
                    cache_epg.filtered_etag, cache_epg.filtered_body
                ),
            )

        headers: dict[str, str] = {}
        if cache_epg.upstream_etag:
            headers["If-None-Match"] = cache_epg.upstream_etag

        resp = await client.get(EPG_URL, headers=headers)
        if resp.status_code == 304:
            if cache_epg.filtered_body:
                cache_epg.last_fetch = now
                if request.headers.get("if-none-match") == cache_epg.filtered_etag:
                    return Response(status_code=304)
                return Response(
                    content=cache_epg.filtered_body,
                    media_type="application/xml",
                    headers=_headers_from_cache(
                        cache_epg.filtered_etag, cache_epg.filtered_body
                    ),
                )
            resp = await client.get(EPG_URL)
        if resp.status_code != 200:
            return Response(status_code=502, content="Upstream EPG unavailable")

        cache_epg.upstream_etag = resp.headers.get("ETag")
        body = filter_epg(resp.content, MTV_TVG_ID, MTV_TVG_NAME)
        cache_epg.filtered_body = body
        cache_epg.filtered_etag = _etag_for(body)
        cache_epg.last_fetch = now

        if request.headers.get("if-none-match") == cache_epg.filtered_etag:
            return Response(status_code=304)

        return Response(
            content=body,
            media_type="application/xml",
            headers=_headers_from_cache(cache_epg.filtered_etag, body),
        )
