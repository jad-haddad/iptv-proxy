package epg

import (
	"encoding/xml"
	"strings"
	"testing"
	"time"
)

func TestGenerateXML_Basic(t *testing.T) {
	loc := time.FixedZone("EEST", 3*3600)
	programmes := []Programme{
		{
			Title:       "Press Review",
			Description: "Daily roundup of headlines.",
			Start:       time.Date(2026, 7, 10, 7, 30, 0, 0, loc),
			Stop:        time.Date(2026, 7, 10, 8, 0, 0, 0, loc),
		},
		{
			Title:       "Morning News",
			Description: "Start your day with news.",
			Start:       time.Date(2026, 7, 10, 8, 0, 0, 0, loc),
			Stop:        time.Date(2026, 7, 10, 8, 30, 0, 0, loc),
		},
	}

	out, err := GenerateXML(programmes, "mtvlebanon.lb", "MTV Lebanon")
	if err != nil {
		t.Fatalf("GenerateXML: %v", err)
	}

	if !strings.HasPrefix(string(out), xml.Header) {
		t.Error("output missing XML header")
	}

	if !strings.Contains(string(out), `<channel id="mtvlebanon.lb">`) {
		t.Error("output missing channel element")
	}
	if !strings.Contains(string(out), `<display-name>MTV Lebanon</display-name>`) {
		t.Error("output missing display-name")
	}

	if !strings.Contains(string(out), `channel="mtvlebanon.lb"`) {
		t.Error("output missing channel attr on programme")
	}

	if !strings.Contains(string(out), `start="20260710073000 +0300"`) {
		t.Error("output missing correct start time")
	}
	if !strings.Contains(string(out), `stop="20260710080000 +0300"`) {
		t.Error("output missing correct stop time")
	}

	if !strings.Contains(string(out), `<title>Press Review</title>`) {
		t.Error("output missing Press Review title")
	}
	if !strings.Contains(string(out), `<desc>Daily roundup of headlines.</desc>`) {
		t.Error("output missing Press Review description")
	}
}

func TestGenerateXML_Empty(t *testing.T) {
	out, err := GenerateXML(nil, "mtvlebanon.lb", "MTV Lebanon")
	if err != nil {
		t.Fatalf("GenerateXML empty: %v", err)
	}

	if !strings.Contains(string(out), xml.Header) {
		t.Error("empty output missing XML header")
	}
	if !strings.Contains(string(out), `<channel id="mtvlebanon.lb">`) {
		t.Error("empty output missing channel")
	}
	if !strings.Contains(string(out), `</tv>`) {
		t.Error("empty output missing closing tv tag")
	}
}

func TestGenerateXML_NoDescription(t *testing.T) {
	loc := time.FixedZone("EET", 2*3600)
	programmes := []Programme{
		{
			Title: "Just a title",
			Start: time.Date(2026, 1, 10, 10, 0, 0, 0, loc),
			Stop:  time.Date(2026, 1, 10, 11, 0, 0, 0, loc),
		},
	}

	out, err := GenerateXML(programmes, "mtvlebanon.lb", "MTV Lebanon")
	if err != nil {
		t.Fatalf("GenerateXML: %v", err)
	}

	if strings.Contains(string(out), `<desc`) {
		t.Error("output should not contain desc element for empty description")
	}
}

func TestGenerateXML_Timestamps(t *testing.T) {
	loc := time.FixedZone("EET", 2*3600)
	programmes := []Programme{
		{
			Title: "Winter Show",
			Start: time.Date(2026, 1, 10, 10, 0, 0, 0, loc),
			Stop:  time.Date(2026, 1, 10, 11, 0, 0, 0, loc),
		},
	}

	out, err := GenerateXML(programmes, "ch", "Channel")
	if err != nil {
		t.Fatalf("GenerateXML: %v", err)
	}

	if !strings.Contains(string(out), `+0200`) {
		t.Errorf("expected +0200 offset for winter: %s", out)
	}
}

func TestGenerateXML_ValidXML(t *testing.T) {
	loc := time.FixedZone("EEST", 3*3600)
	programmes := []Programme{
		{
			Title:       "Show A",
			Description: "Description A",
			Start:       time.Date(2026, 7, 10, 7, 30, 0, 0, loc),
			Stop:        time.Date(2026, 7, 10, 8, 0, 0, 0, loc),
		},
	}

	out, err := GenerateXML(programmes, "mtv", "MTV")
	if err != nil {
		t.Fatalf("GenerateXML: %v", err)
	}

	if err := xml.Unmarshal(out, new(interface{})); err != nil {
		t.Errorf("output is not valid XML: %v", err)
	}
}

func TestGenerateXML_ProgrammeOrder(t *testing.T) {
	loc := time.FixedZone("EEST", 3*3600)
	programmes := []Programme{
		{Title: "First", Start: time.Date(2026, 7, 10, 7, 0, 0, 0, loc), Stop: time.Date(2026, 7, 10, 8, 0, 0, 0, loc)},
		{Title: "Second", Start: time.Date(2026, 7, 10, 8, 0, 0, 0, loc), Stop: time.Date(2026, 7, 10, 9, 0, 0, 0, loc)},
	}

	out, err := GenerateXML(programmes, "ch", "Ch")
	if err != nil {
		t.Fatalf("GenerateXML: %v", err)
	}

	idxFirst := strings.Index(string(out), "<title>First</title>")
	idxSecond := strings.Index(string(out), "<title>Second</title>")

	if idxFirst < 0 || idxSecond < 0 {
		t.Fatal("programme titles not found in output")
	}
	if idxFirst > idxSecond {
		t.Error("programmes should appear in input order; First found after Second")
	}
}

func TestGenerateXML_WithIcon(t *testing.T) {
	loc := time.FixedZone("EEST", 3*3600)
	programmes := []Programme{
		{
			Title: "Show With Icon",
			Icon:  "https://imagescdn.mtv.com.lb/programs/test.jpg",
			Start: time.Date(2026, 7, 10, 7, 30, 0, 0, loc),
			Stop:  time.Date(2026, 7, 10, 8, 0, 0, 0, loc),
		},
	}

	out, err := GenerateXML(programmes, "ch", "Channel")
	if err != nil {
		t.Fatalf("GenerateXML: %v", err)
	}

	if !strings.Contains(string(out), `<icon src="https://imagescdn.mtv.com.lb/programs/test.jpg"`) {
		t.Error("output missing icon element with src attribute")
	}
}

func TestGenerateXML_NoIcon(t *testing.T) {
	loc := time.FixedZone("EEST", 3*3600)
	programmes := []Programme{
		{
			Title: "Show No Icon",
			Start: time.Date(2026, 7, 10, 7, 30, 0, 0, loc),
			Stop:  time.Date(2026, 7, 10, 8, 0, 0, 0, loc),
		},
	}

	out, err := GenerateXML(programmes, "ch", "Channel")
	if err != nil {
		t.Fatalf("GenerateXML: %v", err)
	}

	if strings.Contains(string(out), `<icon`) {
		t.Error("output should not contain icon element for empty Icon")
	}
}

func TestGenerateXML_Rerun(t *testing.T) {
	loc := time.FixedZone("EEST", 3*3600)
	programmes := []Programme{
		{
			Title:   "Evening Rerun",
			IsRerun: true,
			Start:   time.Date(2026, 7, 10, 18, 0, 0, 0, loc),
			Stop:    time.Date(2026, 7, 10, 19, 0, 0, 0, loc),
		},
	}

	out, err := GenerateXML(programmes, "ch", "Channel")
	if err != nil {
		t.Fatalf("GenerateXML: %v", err)
	}

	if !strings.Contains(string(out), `<previously-shown`) {
		t.Error("output missing previously-shown element for rerun programme")
	}
}

func TestGenerateXML_NotRerun(t *testing.T) {
	loc := time.FixedZone("EEST", 3*3600)
	programmes := []Programme{
		{
			Title: "Live Show",
			Start: time.Date(2026, 7, 10, 20, 0, 0, 0, loc),
			Stop:  time.Date(2026, 7, 10, 21, 0, 0, 0, loc),
		},
	}

	out, err := GenerateXML(programmes, "ch", "Channel")
	if err != nil {
		t.Fatalf("GenerateXML: %v", err)
	}

	if strings.Contains(string(out), `<previously-shown`) {
		t.Error("output should not contain previously-shown for non-rerun programme")
	}
}
