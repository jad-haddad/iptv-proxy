package epg

import (
	"encoding/xml"
	"time"
)

type Programme struct {
	Title       string
	Description string
	Icon        string
	IsRerun     bool
	Start       time.Time
	Stop        time.Time
}

type tv struct {
	XMLName    xml.Name      `xml:"tv"`
	Channels   []channel     `xml:"channel"`
	Programmes []tvProgramme `xml:"programme"`
}

type channel struct {
	ID          string `xml:"id,attr"`
	DisplayName string `xml:"display-name"`
}

type icon struct {
	Src string `xml:"src,attr"`
}

type previouslyShown struct{}

type tvProgramme struct {
	Channel         string           `xml:"channel,attr"`
	Start           string           `xml:"start,attr"`
	Stop            string           `xml:"stop,attr"`
	Title           string           `xml:"title"`
	Desc            string           `xml:"desc,omitempty"`
	Icon            *icon            `xml:"icon,omitempty"`
	PreviouslyShown *previouslyShown `xml:"previously-shown,omitempty"`
}

func GenerateXML(programmes []Programme, channelID, channelName string) ([]byte, error) {
	tvProgs := make([]tvProgramme, len(programmes))
	for i, p := range programmes {
		prog := tvProgramme{
			Channel: channelID,
			Start:   p.Start.Format("20060102150405 -0700"),
			Stop:    p.Stop.Format("20060102150405 -0700"),
			Title:   p.Title,
			Desc:    p.Description,
		}
		if p.Icon != "" {
			prog.Icon = &icon{Src: p.Icon}
		}
		if p.IsRerun {
			prog.PreviouslyShown = &previouslyShown{}
		}
		tvProgs[i] = prog
	}

	doc := tv{
		Channels: []channel{
			{ID: channelID, DisplayName: channelName},
		},
		Programmes: tvProgs,
	}

	out, err := xml.MarshalIndent(doc, "", "  ")
	if err != nil {
		return nil, err
	}

	result := make([]byte, 0, len(xml.Header)+len(out))
	result = append(result, []byte(xml.Header)...)
	result = append(result, out...)
	result = append(result, '\n')
	return result, nil
}
