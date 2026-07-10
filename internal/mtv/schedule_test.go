package mtv

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/jad-haddad/iptv-proxy/internal/epg"
)

const daysJSON = `[{"$id":"1","Id":6,"Title":"Monday","Priority":null,"DateCreated":null,"DateModified":null,"IsPublished":false,"IsDeleted":false,"ProgramDailySchedules":[]},{"$id":"2","Id":7,"Title":"Tuesday","Priority":null,"DateCreated":null,"DateModified":null,"IsPublished":false,"IsDeleted":false,"ProgramDailySchedules":[]},{"$id":"3","Id":8,"Title":"Wednesday","Priority":null,"DateCreated":null,"DateModified":null,"IsPublished":false,"IsDeleted":false,"ProgramDailySchedules":[]},{"$id":"4","Id":9,"Title":"Thursday","Priority":null,"DateCreated":null,"DateModified":null,"IsPublished":false,"IsDeleted":false,"ProgramDailySchedules":[]},{"$id":"5","Id":10,"Title":"Friday","Priority":null,"DateCreated":null,"DateModified":null,"IsPublished":false,"IsDeleted":false,"ProgramDailySchedules":[]},{"$id":"6","Id":11,"Title":"Saturday","Priority":null,"DateCreated":null,"DateModified":null,"IsPublished":false,"IsDeleted":false,"ProgramDailySchedules":[]},{"$id":"7","Id":12,"Title":"Sunday","Priority":null,"DateCreated":null,"DateModified":null,"IsPublished":false,"IsDeleted":false,"ProgramDailySchedules":[]}]`

const day10JSON = `[{"$id":"1","Id":91,"Day":10,"ProgramId":26,"SeasonId":2002,"IsPublished":false,"IsDeleted":false,"Time":"07:30     ","Priority":0,"IsRerun":false,"DayNavigation":null,"Program":{"$id":"2","Id":26,"Name":"Press Review","ArName":"Press Review AR","Published":null,"Description":"Daily roundup of headlines.","ArDescription":"AR desc","CreationDate":"0001-01-01T00:00:00","LogoPic":null,"SquarePic":null,"Trailer":null,"CreatedBy":null,"EditedBy":null,"ModifiedDate":null,"CategoryId":null,"StartAiring":null,"StopAiring":null,"MainPic":null,"Slogan":null,"BannerUrl":null,"CreatedByGuid":null,"ModifiedByGuid":null,"DeletedByGuid":null,"VastLink":null,"StopYoutube":false,"IsSpecial":false,"IsArchived":false,"IsFeatured":false,"HashtagName":null,"IsRunning":false,"Deleted":false,"Type":null,"ProgramTypeId":null,"HasNews":false,"HasVideos":false,"HasCast":false,"CastMenuLabel":null,"ArCastMenuLabel":null,"Skin":null,"DateInEpisodesTitle":false,"FacebookId":null,"TwitterId":null,"InstagramId":null,"InstagramUrl":null,"YoutubeId":null,"LandingImgSrc":null,"IsDummy":false,"DummyLink":null,"Dfptag":null,"ShowInReplay":false,"IsTrendingNow":false,"Priority":null,"PictureWitoutTitle":null,"Summary":null,"NewsUrl":null,"Views":null,"Favorites":null,"StationId":null,"AiringTime":null,"PictureWithTitle":null,"PublishDate":null,"UnPublishDate":null,"IsDeleted":false,"Restriction":false,"NewsLink":null,"ArSummary":null,"IsPaid":false,"ShowInWatch":null,"ShowEpisodeDate":null,"IsTenTime":false,"CreatedBy2":null,"ModifiedBy2":null,"DeletedBy2":null,"CustomUrl":null,"CustomPageTitle":null,"CustomH1content":null,"CustomUrlTitle":null,"MetaTitle":null,"MetaTitleSecondary":null,"MetaKeywords":null,"MetaDescription":null,"MetaDescriptionSecondary":null,"MetaDescriptionTertiary":null,"MetaImgSrc":null,"AlternativeNames":null,"ArticleCategoryProgramItems":[],"Articles":[],"Category":null,"CreatedBy2Navigation":null,"DeletedBy2Navigation":null,"FavoriteHistories":[],"ModifiedBy2Navigation":null,"ProgramCategoryItems":[],"ProgramDailySchedules":[],"ProgramDisplays":[],"ProgramGenres":[],"ProgramGeoBlockings":[],"Seasons":[]},"Season":{"$id":"3","Id":2002,"ProgramId":0,"ServerConfigId":null,"Title":"2023","ArTitle":"2023","LogoPic":null,"SquarePic":null,"Trailer":null,"MainPic":"https://example.com/img.jpg","MobileAppImage":null,"HashtagName":null,"PublishStatus":false,"Priority":0,"DateCreated":"0001-01-01T00:00:00","DateModified":"0001-01-01T00:00:00","StartAiring":null,"StopAiring":null,"IsCurrent":false,"Deleted":false,"VideoKey":null,"HasNews":false,"HasVideos":false,"HasCast":false,"LandingImgSrc":null,"VideoKeyImgSrc":null,"Description":null,"ArDescription":null,"AboutImgSrc":null,"IsTenTime":false,"DateUnPublished":null,"PublishDate":null,"UnPublishDate":null,"CreatedBy":null,"ModifiedBy":null,"Articles":[],"Episodes":[],"HomePageSliders":[],"Program":null,"ProgramDailySchedules":[]}},{"$id":"2","Id":799,"Day":10,"ProgramId":28,"SeasonId":1800,"IsPublished":false,"IsDeleted":false,"Time":"08:00     ","Priority":0,"IsRerun":null,"DayNavigation":null,"Program":{"$id":"3","Id":28,"Name":"Morning News","ArName":"Morning News AR","Published":null,"Description":"Start your day with news.","ArDescription":"AR desc","CreationDate":"0001-01-01T00:00:00","LogoPic":null,"SquarePic":null,"Trailer":null,"CreatedBy":null,"EditedBy":null,"ModifiedDate":null,"CategoryId":null,"StartAiring":null,"StopAiring":null,"MainPic":null,"Slogan":null,"BannerUrl":null,"CreatedByGuid":null,"ModifiedByGuid":null,"DeletedByGuid":null,"VastLink":null,"StopYoutube":false,"IsSpecial":false,"IsArchived":false,"IsFeatured":false,"HashtagName":null,"IsRunning":false,"Deleted":false,"Type":null,"ProgramTypeId":null,"HasNews":false,"HasVideos":false,"HasCast":false,"CastMenuLabel":null,"ArCastMenuLabel":null,"Skin":null,"DateInEpisodesTitle":false,"FacebookId":null,"TwitterId":null,"InstagramId":null,"InstagramUrl":null,"YoutubeId":null,"LandingImgSrc":null,"IsDummy":false,"DummyLink":null,"Dfptag":null,"ShowInReplay":false,"IsTrendingNow":false,"Priority":null,"PictureWitoutTitle":null,"Summary":null,"NewsUrl":null,"Views":null,"Favorites":null,"StationId":null,"AiringTime":null,"PictureWithTitle":null,"PublishDate":null,"UnPublishDate":null,"IsDeleted":false,"Restriction":false,"NewsLink":null,"ArSummary":null,"IsPaid":false,"ShowInWatch":null,"ShowEpisodeDate":null,"IsTenTime":false,"CreatedBy2":null,"ModifiedBy2":null,"DeletedBy2":null,"CustomUrl":null,"CustomPageTitle":null,"CustomH1content":null,"CustomUrlTitle":null,"MetaTitle":null,"MetaTitleSecondary":null,"MetaKeywords":null,"MetaDescription":null,"MetaDescriptionSecondary":null,"MetaDescriptionTertiary":null,"MetaImgSrc":null,"AlternativeNames":null,"ArticleCategoryProgramItems":[],"Articles":[],"Category":null,"CreatedBy2Navigation":null,"DeletedBy2Navigation":null,"FavoriteHistories":[],"ModifiedBy2Navigation":null,"ProgramCategoryItems":[],"ProgramDailySchedules":[],"ProgramDisplays":[],"ProgramGenres":[],"ProgramGeoBlockings":[],"Seasons":[]},"Season":{"$id":"4","Id":1800,"ProgramId":0,"ServerConfigId":null,"Title":"2025","ArTitle":"2025","LogoPic":null,"SquarePic":null,"Trailer":null,"MainPic":"https://example.com/img.jpg","MobileAppImage":null,"HashtagName":null,"PublishStatus":false,"Priority":0,"DateCreated":"0001-01-01T00:00:00","DateModified":"0001-01-01T00:00:00","StartAiring":null,"StopAiring":null,"IsCurrent":false,"Deleted":false,"VideoKey":null,"HasNews":false,"HasVideos":false,"HasCast":false,"LandingImgSrc":null,"VideoKeyImgSrc":null,"Description":null,"ArDescription":null,"AboutImgSrc":null,"IsTenTime":false,"DateUnPublished":null,"PublishDate":null,"UnPublishDate":null,"CreatedBy":null,"ModifiedBy":null,"Articles":[],"Episodes":[],"HomePageSliders":[],"Program":null,"ProgramDailySchedules":[]}},{"$id":"3","Id":63,"Day":10,"ProgramId":40,"SeasonId":1008,"IsPublished":false,"IsDeleted":false,"Time":"20:00     ","Priority":0,"IsRerun":null,"DayNavigation":null,"Program":{"$id":"5","Id":40,"Name":"Prime Time News","ArName":"Prime Time News AR","Published":null,"Description":"Evening bulletin.","ArDescription":"AR desc","CreationDate":"0001-01-01T00:00:00","LogoPic":null,"SquarePic":null,"Trailer":null,"CreatedBy":null,"EditedBy":null,"ModifiedDate":null,"CategoryId":null,"StartAiring":null,"StopAiring":null,"MainPic":null,"Slogan":null,"BannerUrl":null,"CreatedByGuid":null,"ModifiedByGuid":null,"DeletedByGuid":null,"VastLink":null,"StopYoutube":false,"IsSpecial":false,"IsArchived":false,"IsFeatured":false,"HashtagName":null,"IsRunning":false,"Deleted":false,"Type":null,"ProgramTypeId":null,"HasNews":false,"HasVideos":false,"HasCast":false,"CastMenuLabel":null,"ArCastMenuLabel":null,"Skin":null,"DateInEpisodesTitle":false,"FacebookId":null,"TwitterId":null,"InstagramId":null,"InstagramUrl":null,"YoutubeId":null,"LandingImgSrc":null,"IsDummy":false,"DummyLink":null,"Dfptag":null,"ShowInReplay":false,"IsTrendingNow":false,"Priority":null,"PictureWitoutTitle":null,"Summary":null,"NewsUrl":null,"Views":null,"Favorites":null,"StationId":null,"AiringTime":null,"PictureWithTitle":null,"PublishDate":null,"UnPublishDate":null,"IsDeleted":false,"Restriction":false,"NewsLink":null,"ArSummary":null,"IsPaid":false,"ShowInWatch":null,"ShowEpisodeDate":null,"IsTenTime":false,"CreatedBy2":null,"ModifiedBy2":null,"DeletedBy2":null,"CustomUrl":null,"CustomPageTitle":null,"CustomH1content":null,"CustomUrlTitle":null,"MetaTitle":null,"MetaTitleSecondary":null,"MetaKeywords":null,"MetaDescription":null,"MetaDescriptionSecondary":null,"MetaDescriptionTertiary":null,"MetaImgSrc":null,"AlternativeNames":null,"ArticleCategoryProgramItems":[],"Articles":[],"Category":null,"CreatedBy2Navigation":null,"DeletedBy2Navigation":null,"FavoriteHistories":[],"ModifiedBy2Navigation":null,"ProgramCategoryItems":[],"ProgramDailySchedules":[],"ProgramDisplays":[],"ProgramGenres":[],"ProgramGeoBlockings":[],"Seasons":[]},"Season":{"$id":"6","Id":1008,"ProgramId":0,"ServerConfigId":null,"Title":"2025","ArTitle":"2025","LogoPic":null,"SquarePic":null,"Trailer":null,"MainPic":"https://example.com/img.jpg","MobileAppImage":null,"HashtagName":null,"PublishStatus":false,"Priority":0,"DateCreated":"0001-01-01T00:00:00","DateModified":"0001-01-01T00:00:00","StartAiring":null,"StopAiring":null,"IsCurrent":false,"Deleted":false,"VideoKey":null,"HasNews":false,"HasVideos":false,"HasCast":false,"LandingImgSrc":null,"VideoKeyImgSrc":null,"Description":null,"ArDescription":null,"AboutImgSrc":null,"IsTenTime":false,"DateUnPublished":null,"PublishDate":null,"UnPublishDate":null,"CreatedBy":null,"ModifiedBy":null,"Articles":[],"Episodes":[],"HomePageSliders":[],"Program":null,"ProgramDailySchedules":[]}}]`

func newTestServer() *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/api/schedule/days":
			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte(daysJSON))
		case "/api/schedule/days/10":
			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte(day10JSON))
		default:
			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte(`[]`))
		}
	}))
}

func TestFetchSchedule(t *testing.T) {
	ts := newTestServer()
	defer ts.Close()

	client := ts.Client()
	friday := time.Date(2026, 7, 10, 12, 0, 0, 0, time.UTC)

	programmes, err := FetchSchedule(client, ts.URL, friday)
	if err != nil {
		t.Fatalf("FetchSchedule: %v", err)
	}

	if len(programmes) < 3 {
		t.Fatalf("expected at least 3 programmes, got %d", len(programmes))
	}

	if programmes[0].Title != "Press Review" {
		t.Errorf("first programme title: got %q, want %q", programmes[0].Title, "Press Review")
	}
	if programmes[0].Description != "Daily roundup of headlines." {
		t.Errorf("first programme desc: got %q", programmes[0].Description)
	}

	if programmes[1].Title != "Morning News" {
		t.Errorf("second programme title: got %q, want %q", programmes[1].Title, "Morning News")
	}

	if programmes[2].Title != "Prime Time News" {
		t.Errorf("third programme title: got %q, want %q", programmes[2].Title, "Prime Time News")
	}

	if !programmes[0].Start.Before(programmes[1].Start) {
		t.Error("programmes not sorted by start time")
	}

	if !programmes[0].Stop.Equal(programmes[1].Start) {
		t.Errorf("first programme stop should equal second programme start: %v != %v", programmes[0].Stop, programmes[1].Start)
	}
}

func TestFetchScheduleServesProgramme(t *testing.T) {
	ts := newTestServer()
	defer ts.Close()

	client := ts.Client()
	friday := time.Date(2026, 7, 10, 12, 0, 0, 0, time.UTC)

	programmes, err := FetchSchedule(client, ts.URL, friday)
	if err != nil {
		t.Fatalf("FetchSchedule: %v", err)
	}

	p := programmes[0]
	if p.Title != "Press Review" {
		t.Errorf("title: got %q, want %q", p.Title, "Press Review")
	}
	if p.Description != "Daily roundup of headlines." {
		t.Errorf("desc: got %q", p.Description)
	}

	startStr := p.Start.Format("20060102150405 -0700")
	stopStr := p.Stop.Format("20060102150405 -0700")

	_ = startStr
	_ = stopStr

	_, err = epg.GenerateXML(programmes, "mtvlebanon.lb", "MTV Lebanon")
	if err != nil {
		t.Fatalf("GenerateXML: %v", err)
	}
}

func TestParseTime(t *testing.T) {
	tests := []struct {
		input    string
		hour     int
		minute   int
		hasError bool
	}{
		{"07:30     ", 7, 30, false},
		{"00:00", 0, 0, false},
		{"23:59", 23, 59, false},
		{"invalid", 0, 0, true},
		{"07:30:00", 0, 0, true},
	}

	for _, tc := range tests {
		parsed, err := parseTime(tc.input)
		if tc.hasError {
			if err == nil {
				t.Errorf("parseTime(%q): expected error, got nil", tc.input)
			}
			continue
		}
		if err != nil {
			t.Errorf("parseTime(%q): %v", tc.input, err)
			continue
		}
		if parsed.Hour() != tc.hour {
			t.Errorf("parseTime(%q): hour %d, want %d", tc.input, parsed.Hour(), tc.hour)
		}
		if parsed.Minute() != tc.minute {
			t.Errorf("parseTime(%q): minute %d, want %d", tc.input, parsed.Minute(), tc.minute)
		}
	}
}

func TestBeirutOffset(t *testing.T) {
	july := time.Date(2026, 7, 10, 12, 0, 0, 0, time.UTC)
	off := beirutOffset(july)
	if off != 3*3600 {
		t.Errorf("July offset: got %d, want %d (EEST UTC+3)", off, 3*3600)
	}

	jan := time.Date(2026, 1, 10, 12, 0, 0, 0, time.UTC)
	off = beirutOffset(jan)
	if off != 2*3600 {
		t.Errorf("January offset: got %d, want %d (EET UTC+2)", off, 2*3600)
	}
}

func TestDayToWeekday(t *testing.T) {
	tests := map[int]time.Weekday{
		6:  time.Monday,
		7:  time.Tuesday,
		8:  time.Wednesday,
		9:  time.Thursday,
		10: time.Friday,
		11: time.Saturday,
		12: time.Sunday,
	}

	for id, expected := range tests {
		got, ok := dayToWeekday[id]
		if !ok {
			t.Errorf("dayToWeekday[%d]: not found", id)
			continue
		}
		if got != expected {
			t.Errorf("dayToWeekday[%d]: got %v, want %v", id, got, expected)
		}
	}
}

func TestComputeDayDates(t *testing.T) {
	friday := time.Date(2026, 7, 10, 12, 0, 0, 0, time.UTC).In(time.UTC)
	days := daysResponse{
		{ID: 10, Title: "Friday"},
		{ID: 11, Title: "Saturday"},
		{ID: 6, Title: "Monday"},
	}

	dates := computeDayDates(friday, days)

	fridayDate := dates[10]
	if fridayDate.Year() != 2026 || fridayDate.Month() != 7 || fridayDate.Day() != 10 {
		t.Errorf("Friday date: got %s, want 2026-07-10", fridayDate.Format("2006-01-02"))
	}

	saturdayDate := dates[11]
	if saturdayDate.Year() != 2026 || saturdayDate.Month() != 7 || saturdayDate.Day() != 11 {
		t.Errorf("Saturday date: got %s, want 2026-07-11", saturdayDate.Format("2006-01-02"))
	}

	mondayDate := dates[6]
	if mondayDate.Year() != 2026 || mondayDate.Month() != 7 || mondayDate.Day() != 13 {
		t.Errorf("Monday date: got %s, want 2026-07-13", mondayDate.Format("2006-01-02"))
	}
}

func TestFetchScheduleServerError(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer ts.Close()

	client := ts.Client()
	_, err := FetchSchedule(client, ts.URL, time.Now())
	if err == nil {
		t.Error("expected error on server 500")
	}
}
