import os


def _get_env(name: str, default: str) -> str:
    value = os.environ.get(name)
    return value if value else default


M3U_URL = _get_env("M3U_URL", "https://iptv-org.github.io/iptv/countries/lb.m3u")
EPG_URL = _get_env("EPG_URL", "https://mdag9904.github.io/lebanon-epg/epg.xml")

MTV_REGEX = _get_env(
    "MTV_REGEX",
    r"(?i)\bmtv\b.*\blebanon\b|mtv\s*lebanon|mtvlebanon",
)
MTV_TVG_ID = _get_env("MTV_TVG_ID", "mtvlebanon.lb")
MTV_TVG_NAME = _get_env("MTV_TVG_NAME", "MTV Lebanon")

EPG_REFRESH_SECONDS = int(_get_env("EPG_REFRESH_SECONDS", "3600"))
REQUEST_TIMEOUT_SECONDS = float(_get_env("REQUEST_TIMEOUT_SECONDS", "15"))
