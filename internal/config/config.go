package config

import (
	"os"
	"strconv"
	"strings"
)

type Config struct {
	M3UURL                string
	EPGURL                string
	MTVRegex              string
	MTVTVGID              string
	MTVTVGName            string
	EPGRefreshSeconds     int
	RequestTimeoutSeconds float64
}

func Load() Config {
	return Config{
		M3UURL:                getEnv("M3U_URL", "https://iptv-org.github.io/iptv/countries/lb.m3u"),
		EPGURL:                getEnv("EPG_URL", "https://mdag9904.github.io/lebanon-epg/epg.xml"),
		MTVRegex:              getEnv("MTV_REGEX", `(?i)\bmtv\b.*\blebanon\b|mtv\s*lebanon|mtvlebanon`),
		MTVTVGID:              getEnv("MTV_TVG_ID", "mtvlebanon.lb"),
		MTVTVGName:            getEnv("MTV_TVG_NAME", "MTV Lebanon"),
		EPGRefreshSeconds:     getEnvInt("EPG_REFRESH_SECONDS", 3600),
		RequestTimeoutSeconds: getEnvFloat("REQUEST_TIMEOUT_SECONDS", 15),
	}
}

func getEnv(name, def string) string {
	value := strings.TrimSpace(os.Getenv(name))
	if value == "" {
		return def
	}
	return value
}

func getEnvInt(name string, def int) int {
	value := strings.TrimSpace(os.Getenv(name))
	if value == "" {
		return def
	}
	parsed, err := strconv.Atoi(value)
	if err != nil {
		return def
	}
	return parsed
}

func getEnvFloat(name string, def float64) float64 {
	value := strings.TrimSpace(os.Getenv(name))
	if value == "" {
		return def
	}
	parsed, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return def
	}
	return parsed
}
