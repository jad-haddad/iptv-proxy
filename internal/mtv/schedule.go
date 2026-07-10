package mtv

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/jad-haddad/iptv-proxy/internal/epg"
)

var dayToWeekday = map[int]time.Weekday{
	6:  time.Monday,
	7:  time.Tuesday,
	8:  time.Wednesday,
	9:  time.Thursday,
	10: time.Friday,
	11: time.Saturday,
	12: time.Sunday,
}

type dayEntry struct {
	date  time.Time
	entry scheduleEntry
}

func (d dayEntry) dateTime() time.Time {
	t, _ := parseTime(d.entry.Time)
	return time.Date(d.date.Year(), d.date.Month(), d.date.Day(), t.Hour(), t.Minute(), 0, 0, beirutLocation(d.date))
}

func FetchSchedule(client *http.Client, baseURL string, now time.Time) ([]epg.Programme, error) {
	days, err := fetchDays(client, baseURL)
	if err != nil {
		return nil, err
	}

	beirutNow := now.In(beirutLocation(now))
	dayDates := computeDayDates(beirutNow, days)

	var allEntries []dayEntry
	for _, d := range days {
		entries, err := fetchDaySchedule(client, baseURL, d.ID)
		if err != nil {
			return nil, err
		}
		date, ok := dayDates[d.ID]
		if !ok {
			continue
		}
		for _, e := range entries {
			allEntries = append(allEntries, dayEntry{date: date, entry: e})
		}
	}

	sort.Slice(allEntries, func(i, j int) bool {
		return allEntries[i].dateTime().Before(allEntries[j].dateTime())
	})

	programmes := make([]epg.Programme, 0, len(allEntries))
	for i, e := range allEntries {
		start := e.dateTime()
		var stop time.Time
		if i+1 < len(allEntries) {
			stop = allEntries[i+1].dateTime()
		} else {
			stop = start.Add(60 * time.Minute)
		}
		programmes = append(programmes, epg.Programme{
			Title:       e.entry.Program.Name,
			Description: e.entry.Program.Description,
			Start:       start,
			Stop:        stop,
		})
	}

	return programmes, nil
}

func fetchDays(client *http.Client, baseURL string) (daysResponse, error) {
	req, err := http.NewRequest(http.MethodGet, baseURL+"/api/schedule/days", nil)
	if err != nil {
		return nil, err
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("mtv days: status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var days daysResponse
	if err := json.Unmarshal(body, &days); err != nil {
		return nil, err
	}

	return days, nil
}

func fetchDaySchedule(client *http.Client, baseURL string, dayID int) (scheduleResponse, error) {
	url := fmt.Sprintf("%s/api/schedule/days/%d", baseURL, dayID)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("mtv schedule day %d: status %d", dayID, resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var entries scheduleResponse
	if err := json.Unmarshal(body, &entries); err != nil {
		return nil, err
	}

	return entries, nil
}

func computeDayDates(now time.Time, days daysResponse) map[int]time.Time {
	today := now.Weekday()
	result := make(map[int]time.Time, len(days))
	for _, d := range days {
		target, ok := dayToWeekday[d.ID]
		if !ok {
			continue
		}
		daysAhead := (int(target) - int(today) + 7) % 7
		date := now.AddDate(0, 0, daysAhead)
		result[d.ID] = time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
	}
	return result
}

func parseTime(t string) (time.Time, error) {
	trimmed := strings.TrimSpace(t)
	parts := strings.SplitN(trimmed, ":", 2)
	if len(parts) != 2 {
		return time.Time{}, fmt.Errorf("invalid time: %q", t)
	}
	hour, err := strconv.Atoi(parts[0])
	if err != nil {
		return time.Time{}, fmt.Errorf("invalid hour: %q", t)
	}
	minute, err := strconv.Atoi(parts[1])
	if err != nil {
		return time.Time{}, fmt.Errorf("invalid minute: %q", t)
	}
	return time.Date(0, 1, 1, hour, minute, 0, 0, time.UTC), nil
}

func beirutLocation(t time.Time) *time.Location {
	offset := beirutOffset(t)
	return time.FixedZone("EET", offset)
}

func beirutOffset(t time.Time) int {
	year := t.Year()

	marLast := time.Date(year, 3, 31, 0, 0, 0, 0, time.UTC)
	for marLast.Weekday() != time.Sunday {
		marLast = marLast.AddDate(0, 0, -1)
	}

	octLast := time.Date(year, 10, 31, 0, 0, 0, 0, time.UTC)
	for octLast.Weekday() != time.Sunday {
		octLast = octLast.AddDate(0, 0, -1)
	}

	marStart := time.Date(year, 3, marLast.Day(), 0, 0, 0, 0, time.UTC)
	octStart := time.Date(year, 10, octLast.Day(), 0, 0, 0, 0, time.UTC)

	if !t.Before(marStart) && t.Before(octStart) {
		return 3 * 3600
	}
	return 2 * 3600
}
