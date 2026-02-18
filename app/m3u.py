import re
from typing import Iterable, List, Tuple


ATTR_RE = re.compile(r'([\w-]+)="([^"]*)"')


def _parse_extinf(line: str) -> Tuple[List[Tuple[str, str]], str]:
    if "," in line:
        meta, title = line.split(",", 1)
    else:
        meta, title = line, ""
    attrs: List[Tuple[str, str]] = []
    for key, value in ATTR_RE.findall(meta):
        attrs.append((key, value))
    return attrs, title.strip()


def _update_attrs(
    attrs: List[Tuple[str, str]],
    tvg_id: str,
    tvg_name: str,
) -> List[Tuple[str, str]]:
    updated: List[Tuple[str, str]] = []
    seen = {"tvg-id": False, "tvg-name": False}
    for key, value in attrs:
        if key == "tvg-id":
            updated.append((key, tvg_id))
            seen["tvg-id"] = True
        elif key == "tvg-name":
            updated.append((key, tvg_name))
            seen["tvg-name"] = True
        else:
            updated.append((key, value))
    if not seen["tvg-id"]:
        updated.insert(0, ("tvg-id", tvg_id))
    if not seen["tvg-name"]:
        insert_at = 1 if updated and updated[0][0] == "tvg-id" else 0
        updated.insert(insert_at, ("tvg-name", tvg_name))
    return updated


def _render_extinf(attrs: List[Tuple[str, str]], title: str) -> str:
    parts = ["#EXTINF:-1"]
    for key, value in attrs:
        parts.append(f'{key}="{value}"')
    return f"{' '.join(parts)},{title}".strip()


def _match_texts(pattern: re.Pattern[str], texts: Iterable[str | None]) -> bool:
    for text in texts:
        if text and pattern.search(text):
            return True
    return False


def filter_m3u(
    raw: str,
    regex: str,
    tvg_id: str,
    tvg_name: str,
) -> str:
    pattern = re.compile(regex)
    lines = raw.splitlines()
    output: List[str] = ["#EXTM3U"]

    current_meta: List[str] = []
    current_extinf_attrs: List[Tuple[str, str]] = []
    current_title = ""
    matched = False

    for line in lines:
        stripped = line.strip()
        if not stripped:
            continue
        if stripped.startswith("#EXTM3U"):
            continue
        if stripped.startswith("#EXTINF"):
            current_meta = [stripped]
            current_extinf_attrs, current_title = _parse_extinf(stripped)
            matched = _match_texts(
                pattern,
                [
                    dict(current_extinf_attrs).get("tvg-id"),
                    dict(current_extinf_attrs).get("tvg-name"),
                    current_title,
                ],
            )
            continue
        if stripped.startswith("#"):
            if current_meta:
                current_meta.append(stripped)
            continue

        url = stripped
        if current_meta and matched:
            updated_attrs = _update_attrs(current_extinf_attrs, tvg_id, tvg_name)
            output.append(_render_extinf(updated_attrs, tvg_name))
            for meta in current_meta[1:]:
                output.append(meta)
            output.append(url)
            output.append("")
        current_meta = []
        current_extinf_attrs = []
        current_title = ""
        matched = False

    return "\n".join(output).strip() + "\n"
