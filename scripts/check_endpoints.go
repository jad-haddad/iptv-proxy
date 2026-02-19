package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

func main() {
	baseURL := "http://127.0.0.1:8080"
	if len(os.Args) > 1 {
		baseURL = strings.TrimRight(os.Args[1], "/")
	}

	m3uURL := baseURL + "/lebanon.m3u"
	epgURL := baseURL + "/epg.xml"

	client := &http.Client{Timeout: 10 * time.Second}

	m3u, err := client.Get(m3uURL)
	if err != nil {
		fmt.Printf("M3U request failed: %v\n", err)
		os.Exit(1)
	}
	m3uBody, _ := io.ReadAll(m3u.Body)
	_ = m3u.Body.Close()
	if m3u.StatusCode != http.StatusOK {
		fmt.Printf("M3U request failed: %d\n", m3u.StatusCode)
		os.Exit(1)
	}
	if !strings.Contains(string(m3uBody), "#EXTM3U") {
		fmt.Println("M3U response missing #EXTM3U")
		os.Exit(1)
	}
	if !strings.Contains(string(m3uBody), `tvg-id="mtvlebanon.lb"`) {
		fmt.Println("M3U response missing normalized tvg-id")
		os.Exit(1)
	}

	m3uETag := m3u.Header.Get("ETag")
	if m3uETag == "" {
		fmt.Println("M3U response missing ETag")
		os.Exit(1)
	}
	m3uReq, err := http.NewRequest(http.MethodGet, m3uURL, nil)
	if err != nil {
		fmt.Printf("M3U ETag cache check failed: %v\n", err)
		os.Exit(1)
	}
	m3uReq.Header.Set("If-None-Match", m3uETag)
	m3uCached, err := client.Do(m3uReq)
	if err != nil {
		fmt.Printf("M3U ETag cache check failed: %v\n", err)
		os.Exit(1)
	}
	m3uCached.Body.Close()
	if m3uCached.StatusCode != http.StatusNotModified {
		fmt.Printf("M3U ETag cache check failed: %d\n", m3uCached.StatusCode)
		os.Exit(1)
	}

	epg, err := client.Get(epgURL)
	if err != nil {
		fmt.Printf("EPG request failed: %v\n", err)
		os.Exit(1)
	}
	epgBody, _ := io.ReadAll(epg.Body)
	_ = epg.Body.Close()
	if epg.StatusCode != http.StatusOK {
		fmt.Printf("EPG request failed: %d\n", epg.StatusCode)
		os.Exit(1)
	}
	if !strings.Contains(string(epgBody), `<channel id="mtvlebanon.lb">`) {
		fmt.Println("EPG response missing MTV channel")
		os.Exit(1)
	}
	if !strings.Contains(string(epgBody), `channel="mtvlebanon.lb"`) {
		fmt.Println("EPG response missing programme entries")
		os.Exit(1)
	}

	epgETag := epg.Header.Get("ETag")
	if epgETag == "" {
		fmt.Println("EPG response missing ETag")
		os.Exit(1)
	}
	epgReq, err := http.NewRequest(http.MethodGet, epgURL, nil)
	if err != nil {
		fmt.Printf("EPG ETag cache check failed: %v\n", err)
		os.Exit(1)
	}
	epgReq.Header.Set("If-None-Match", epgETag)
	epgCached, err := client.Do(epgReq)
	if err != nil {
		fmt.Printf("EPG ETag cache check failed: %v\n", err)
		os.Exit(1)
	}
	epgCached.Body.Close()
	if epgCached.StatusCode != http.StatusNotModified {
		fmt.Printf("EPG ETag cache check failed: %d\n", epgCached.StatusCode)
		os.Exit(1)
	}

	fmt.Println("OK")
}
