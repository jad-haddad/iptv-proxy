package m3u

import (
	"fmt"
	"regexp"
	"strings"
)

type AttrPair struct {
	Key   string
	Value string
}

var attrRegex = regexp.MustCompile(`([\w-]+)="([^"]*)"`)

func Filter(raw string, pattern *regexp.Regexp, tvgID, tvgName string) string {
	lines := strings.Split(raw, "\n")
	output := []string{"#EXTM3U"}

	currentMeta := make([]string, 0)
	currentAttrs := make([]AttrPair, 0)
	currentTitle := ""
	matched := false

	for _, line := range lines {
		stripped := strings.TrimSpace(line)
		if stripped == "" {
			continue
		}
		if strings.HasPrefix(stripped, "#EXTM3U") {
			continue
		}
		if strings.HasPrefix(stripped, "#EXTINF") {
			currentMeta = []string{stripped}
			currentAttrs, currentTitle = parseExtinf(stripped)
			attrsMap := make(map[string]string, len(currentAttrs))
			for _, attr := range currentAttrs {
				attrsMap[attr.Key] = attr.Value
			}
			matched = matchTexts(pattern, []string{attrsMap["tvg-id"], attrsMap["tvg-name"], currentTitle})
			continue
		}
		if strings.HasPrefix(stripped, "#") {
			if len(currentMeta) > 0 {
				currentMeta = append(currentMeta, stripped)
			}
			continue
		}

		url := stripped
		if len(currentMeta) > 0 && matched {
			updatedAttrs := updateAttrs(currentAttrs, tvgID, tvgName)
			output = append(output, renderExtinf(updatedAttrs, tvgName))
			if len(currentMeta) > 1 {
				output = append(output, currentMeta[1:]...)
			}
			output = append(output, url, "")
		}
		currentMeta = nil
		currentAttrs = nil
		currentTitle = ""
		matched = false
	}

	return strings.TrimSpace(strings.Join(output, "\n")) + "\n"
}

func parseExtinf(line string) ([]AttrPair, string) {
	meta := line
	title := ""
	if idx := strings.Index(line, ","); idx >= 0 {
		meta = line[:idx]
		title = line[idx+1:]
	}

	attrs := make([]AttrPair, 0)
	for _, match := range attrRegex.FindAllStringSubmatch(meta, -1) {
		attrs = append(attrs, AttrPair{Key: match[1], Value: match[2]})
	}
	return attrs, strings.TrimSpace(title)
}

func updateAttrs(attrs []AttrPair, tvgID, tvgName string) []AttrPair {
	updated := make([]AttrPair, 0, len(attrs)+2)
	seenID := false
	seenName := false
	for _, attr := range attrs {
		switch attr.Key {
		case "tvg-id":
			updated = append(updated, AttrPair{Key: attr.Key, Value: tvgID})
			seenID = true
		case "tvg-name":
			updated = append(updated, AttrPair{Key: attr.Key, Value: tvgName})
			seenName = true
		default:
			updated = append(updated, attr)
		}
	}
	if !seenID {
		updated = append([]AttrPair{{Key: "tvg-id", Value: tvgID}}, updated...)
	}
	if !seenName {
		insertAt := 0
		if len(updated) > 0 && updated[0].Key == "tvg-id" {
			insertAt = 1
		}
		updated = append(updated[:insertAt], append([]AttrPair{{Key: "tvg-name", Value: tvgName}}, updated[insertAt:]...)...)
	}
	return updated
}

func renderExtinf(attrs []AttrPair, title string) string {
	parts := []string{"#EXTINF:-1"}
	for _, attr := range attrs {
		parts = append(parts, fmt.Sprintf("%s=\"%s\"", attr.Key, attr.Value))
	}
	return strings.TrimSpace(strings.Join(parts, " ") + "," + title)
}

func matchTexts(pattern *regexp.Regexp, texts []string) bool {
	for _, text := range texts {
		if text != "" && pattern.MatchString(text) {
			return true
		}
	}
	return false
}
