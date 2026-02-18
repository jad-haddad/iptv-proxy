import asyncio
from dataclasses import dataclass, field


@dataclass
class M3UCache:
    upstream_etag: str | None = None
    filtered_etag: str | None = None
    filtered_body: bytes | None = None
    lock: asyncio.Lock = field(default_factory=asyncio.Lock)


@dataclass
class EPGCache:
    upstream_etag: str | None = None
    filtered_etag: str | None = None
    filtered_body: bytes | None = None
    last_fetch: float = 0.0
    lock: asyncio.Lock = field(default_factory=asyncio.Lock)
