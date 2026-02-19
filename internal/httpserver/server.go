package httpserver

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"time"

	"github.com/jad-haddad/iptv-proxy/internal/cache"
	"github.com/jad-haddad/iptv-proxy/internal/config"
	"github.com/jad-haddad/iptv-proxy/internal/epg"
	"github.com/jad-haddad/iptv-proxy/internal/m3u"
)

type Server struct {
	cfg        config.Config
	client     *http.Client
	m3uCache   cache.M3UCache
	epgCache   cache.EPGCache
	m3uPattern *regexp.Regexp
}

func New(cfg config.Config) (*http.Server, error) {
	pattern, err := regexp.Compile(cfg.MTVRegex)
	if err != nil {
		return nil, err
	}

	svc := &Server{
		cfg:        cfg,
		client:     &http.Client{Timeout: time.Duration(cfg.RequestTimeoutSeconds * float64(time.Second))},
		m3uPattern: pattern,
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/health", svc.health)
	mux.HandleFunc("/lebanon.m3u", svc.getM3U)
	mux.HandleFunc("/epg.xml", svc.getEPG)

	return &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}, nil
}

func (s *Server) health(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(`{"status":"ok"}`))
}

func (s *Server) getM3U(w http.ResponseWriter, r *http.Request) {
	if s.m3uCache.FilteredBody != nil && r.Header.Get("If-None-Match") == s.m3uCache.FilteredETag {
		w.WriteHeader(http.StatusNotModified)
		return
	}

	s.m3uCache.Mu.Lock()
	defer s.m3uCache.Mu.Unlock()

	if s.m3uCache.FilteredBody != nil && r.Header.Get("If-None-Match") == s.m3uCache.FilteredETag {
		w.WriteHeader(http.StatusNotModified)
		return
	}

	req, err := http.NewRequest(http.MethodGet, s.cfg.M3UURL, nil)
	if err != nil {
		http.Error(w, "Upstream M3U unavailable", http.StatusBadGateway)
		return
	}
	if s.m3uCache.UpstreamETag != "" {
		req.Header.Set("If-None-Match", s.m3uCache.UpstreamETag)
	}

	resp, err := s.client.Do(req)
	if err != nil {
		http.Error(w, "Upstream M3U unavailable", http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotModified {
		if s.m3uCache.FilteredBody != nil {
			w.Header().Set("Content-Type", "application/vnd.apple.mpegurl")
			setCacheHeaders(w, s.m3uCache.FilteredETag)
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write(s.m3uCache.FilteredBody)
			return
		}
		req, err = http.NewRequest(http.MethodGet, s.cfg.M3UURL, nil)
		if err != nil {
			http.Error(w, "Upstream M3U unavailable", http.StatusBadGateway)
			return
		}
		resp, err = s.client.Do(req)
		if err != nil {
			http.Error(w, "Upstream M3U unavailable", http.StatusBadGateway)
			return
		}
		defer resp.Body.Close()
	}

	if resp.StatusCode != http.StatusOK {
		http.Error(w, "Upstream M3U unavailable", http.StatusBadGateway)
		return
	}

	raw, err := io.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, "Upstream M3U unavailable", http.StatusBadGateway)
		return
	}

	s.m3uCache.UpstreamETag = resp.Header.Get("ETag")
	filtered := m3u.Filter(string(raw), s.m3uPattern, s.cfg.MTVTVGID, s.cfg.MTVTVGName)
	body := []byte(filtered)
	s.m3uCache.FilteredBody = body
	s.m3uCache.FilteredETag = etagFor(body)

	if r.Header.Get("If-None-Match") == s.m3uCache.FilteredETag {
		w.WriteHeader(http.StatusNotModified)
		return
	}

	w.Header().Set("Content-Type", "application/vnd.apple.mpegurl")
	setCacheHeaders(w, s.m3uCache.FilteredETag)
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(body)
}

func (s *Server) getEPG(w http.ResponseWriter, r *http.Request) {
	now := time.Now()
	if s.epgCache.FilteredBody != nil && now.Sub(s.epgCache.LastFetch) < time.Duration(s.cfg.EPGRefreshSeconds)*time.Second && r.Header.Get("If-None-Match") == s.epgCache.FilteredETag {
		w.WriteHeader(http.StatusNotModified)
		return
	}

	s.epgCache.Mu.Lock()
	defer s.epgCache.Mu.Unlock()

	now = time.Now()
	if s.epgCache.FilteredBody != nil && now.Sub(s.epgCache.LastFetch) < time.Duration(s.cfg.EPGRefreshSeconds)*time.Second {
		if r.Header.Get("If-None-Match") == s.epgCache.FilteredETag {
			w.WriteHeader(http.StatusNotModified)
			return
		}
		w.Header().Set("Content-Type", "application/xml")
		setCacheHeaders(w, s.epgCache.FilteredETag)
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(s.epgCache.FilteredBody)
		return
	}

	req, err := http.NewRequest(http.MethodGet, s.cfg.EPGURL, nil)
	if err != nil {
		http.Error(w, "Upstream EPG unavailable", http.StatusBadGateway)
		return
	}
	if s.epgCache.UpstreamETag != "" {
		req.Header.Set("If-None-Match", s.epgCache.UpstreamETag)
	}

	resp, err := s.client.Do(req)
	if err != nil {
		http.Error(w, "Upstream EPG unavailable", http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotModified {
		if s.epgCache.FilteredBody != nil {
			s.epgCache.LastFetch = now
			if r.Header.Get("If-None-Match") == s.epgCache.FilteredETag {
				w.WriteHeader(http.StatusNotModified)
				return
			}
			w.Header().Set("Content-Type", "application/xml")
			setCacheHeaders(w, s.epgCache.FilteredETag)
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write(s.epgCache.FilteredBody)
			return
		}
		req, err = http.NewRequest(http.MethodGet, s.cfg.EPGURL, nil)
		if err != nil {
			http.Error(w, "Upstream EPG unavailable", http.StatusBadGateway)
			return
		}
		resp, err = s.client.Do(req)
		if err != nil {
			http.Error(w, "Upstream EPG unavailable", http.StatusBadGateway)
			return
		}
		defer resp.Body.Close()
	}

	if resp.StatusCode != http.StatusOK {
		http.Error(w, "Upstream EPG unavailable", http.StatusBadGateway)
		return
	}

	raw, err := io.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, "Upstream EPG unavailable", http.StatusBadGateway)
		return
	}

	s.epgCache.UpstreamETag = resp.Header.Get("ETag")
	filtered, err := epg.Filter(raw, s.cfg.MTVTVGID, s.cfg.MTVTVGName)
	if err != nil {
		http.Error(w, "Upstream EPG unavailable", http.StatusBadGateway)
		return
	}
	s.epgCache.FilteredBody = filtered
	s.epgCache.FilteredETag = etagFor(filtered)
	s.epgCache.LastFetch = now

	if r.Header.Get("If-None-Match") == s.epgCache.FilteredETag {
		w.WriteHeader(http.StatusNotModified)
		return
	}

	w.Header().Set("Content-Type", "application/xml")
	setCacheHeaders(w, s.epgCache.FilteredETag)
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(filtered)
}

func etagFor(body []byte) string {
	sum := sha256.Sum256(body)
	return fmt.Sprintf("\"%s\"", hex.EncodeToString(sum[:]))
}

func setCacheHeaders(w http.ResponseWriter, etag string) {
	if etag == "" {
		return
	}
	w.Header().Set("ETag", etag)
	w.Header().Set("Cache-Control", "no-cache")
}
