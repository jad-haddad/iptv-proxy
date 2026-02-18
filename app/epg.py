import xml.etree.ElementTree as ET


def filter_epg(raw: bytes, channel_id: str, channel_name: str) -> bytes:
    root = ET.fromstring(raw)
    new_root = ET.Element(root.tag, root.attrib)

    channel_found = False
    for channel in root.findall("channel"):
        if channel.get("id") == channel_id:
            new_root.append(channel)
            channel_found = True
            break

    if not channel_found:
        channel = ET.Element("channel", {"id": channel_id})
        display_name = ET.SubElement(channel, "display-name")
        display_name.text = channel_name
        new_root.append(channel)

    for programme in root.findall("programme"):
        if programme.get("channel") == channel_id:
            new_root.append(programme)

    return ET.tostring(new_root, encoding="utf-8", xml_declaration=True)
