import sys

import httpx


def main() -> int:
    base_url = sys.argv[1] if len(sys.argv) > 1 else "http://127.0.0.1:8080"
    m3u_url = f"{base_url.rstrip('/')}/lebanon.m3u"
    epg_url = f"{base_url.rstrip('/')}/epg.xml"

    with httpx.Client(timeout=10.0) as client:
        m3u = client.get(m3u_url)
        if m3u.status_code != 200:
            print(f"M3U request failed: {m3u.status_code}")
            return 1
        if "#EXTM3U" not in m3u.text:
            print("M3U response missing #EXTM3U")
            return 1
        if 'tvg-id="mtvlebanon.lb"' not in m3u.text:
            print("M3U response missing normalized tvg-id")
            return 1

        m3u_etag = m3u.headers.get("ETag")
        if not m3u_etag:
            print("M3U response missing ETag")
            return 1
        m3u_cached = client.get(m3u_url, headers={"If-None-Match": m3u_etag})
        if m3u_cached.status_code != 304:
            print(f"M3U ETag cache check failed: {m3u_cached.status_code}")
            return 1

        epg = client.get(epg_url)
        if epg.status_code != 200:
            print(f"EPG request failed: {epg.status_code}")
            return 1
        if '<channel id="mtvlebanon.lb">' not in epg.text:
            print("EPG response missing MTV channel")
            return 1
        if 'channel="mtvlebanon.lb"' not in epg.text:
            print("EPG response missing programme entries")
            return 1

        epg_etag = epg.headers.get("ETag")
        if not epg_etag:
            print("EPG response missing ETag")
            return 1
        epg_cached = client.get(epg_url, headers={"If-None-Match": epg_etag})
        if epg_cached.status_code != 304:
            print(f"EPG ETag cache check failed: {epg_cached.status_code}")
            return 1

    print("OK")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
