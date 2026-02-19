package epg

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io"
)

type RawElement struct {
	XMLName  xml.Name
	Attr     []xml.Attr `xml:",any,attr"`
	InnerXML string     `xml:",innerxml"`
}

func Filter(raw []byte, channelID, channelName string) ([]byte, error) {
	decoder := xml.NewDecoder(bytes.NewReader(raw))
	var root xml.StartElement
	channels := make([]RawElement, 0)
	programmes := make([]RawElement, 0)
	channelFound := false
	depth := 0

	for {
		tok, err := decoder.Token()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		switch t := tok.(type) {
		case xml.StartElement:
			if depth == 0 {
				root = t
				depth = 1
				continue
			}
			if depth == 1 && t.Name.Local == "channel" {
				var el RawElement
				if err := decoder.DecodeElement(&el, &t); err != nil {
					return nil, err
				}
				if !channelFound && attrValue(el.Attr, "id") == channelID {
					channels = append(channels, el)
					channelFound = true
				}
				continue
			}
			if depth == 1 && t.Name.Local == "programme" {
				var el RawElement
				if err := decoder.DecodeElement(&el, &t); err != nil {
					return nil, err
				}
				if attrValue(el.Attr, "channel") == channelID {
					programmes = append(programmes, el)
				}
				continue
			}
			depth++
		case xml.EndElement:
			depth--
		}
	}

	if root.Name.Local == "" {
		return nil, fmt.Errorf("missing root element")
	}

	var buf bytes.Buffer
	buf.WriteString(xml.Header)
	enc := xml.NewEncoder(&buf)

	if err := enc.EncodeToken(root); err != nil {
		return nil, err
	}

	if channelFound {
		for _, ch := range channels {
			if err := writeRawElement(enc, &buf, ch); err != nil {
				return nil, err
			}
		}
	} else {
		chStart := xml.StartElement{
			Name: xml.Name{Local: "channel"},
			Attr: []xml.Attr{{Name: xml.Name{Local: "id"}, Value: channelID}},
		}
		if err := enc.EncodeToken(chStart); err != nil {
			return nil, err
		}
		if err := enc.EncodeElement(channelName, xml.StartElement{Name: xml.Name{Local: "display-name"}}); err != nil {
			return nil, err
		}
		if err := enc.EncodeToken(xml.EndElement{Name: chStart.Name}); err != nil {
			return nil, err
		}
	}

	for _, pr := range programmes {
		if err := writeRawElement(enc, &buf, pr); err != nil {
			return nil, err
		}
	}

	if err := enc.EncodeToken(xml.EndElement{Name: root.Name}); err != nil {
		return nil, err
	}
	if err := enc.Flush(); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func attrValue(attrs []xml.Attr, name string) string {
	for _, attr := range attrs {
		if attr.Name.Local == name {
			return attr.Value
		}
	}
	return ""
}

func writeRawElement(enc *xml.Encoder, buf *bytes.Buffer, el RawElement) error {
	start := xml.StartElement{Name: el.XMLName, Attr: el.Attr}
	if err := enc.EncodeToken(start); err != nil {
		return err
	}
	if el.InnerXML != "" {
		if err := enc.Flush(); err != nil {
			return err
		}
		buf.WriteString(el.InnerXML)
	}
	if err := enc.EncodeToken(xml.EndElement{Name: el.XMLName}); err != nil {
		return err
	}
	return nil
}
