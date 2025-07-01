package marprom

import (
	"bytes"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"log"
	"regexp"
	"strconv"
	"strings"
)

var stopIDRegex = regexp.MustCompile(`stop=(\d+)`)

type HTMLParser struct{}

func NewHTMLParser() *HTMLParser {
	return &HTMLParser{}
}

func (p *HTMLParser) ParseBusStations(html []byte) ([]BusStation, error) {
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(html))
	if err != nil {
		err := fmt.Errorf("failed to parse HTML: %w", err)
		log.Fatalf("Error parsing HTML: %s", err)
		return nil, err
	}

	locations := extractLocations(doc)

	var stations []BusStation
	doc.Find("#TableOfStops tr").Each(func(_ int, tr *goquery.Selection) {
		if station := parseStationRow(tr, locations); station != nil {
			stations = append(stations, *station)
		}
	})

	return stations, nil
}

func extractLocations(doc *goquery.Document) map[string][2]float64 {
	locations := make(map[string][2]float64)

	doc.Find("#stopStopPoint option").Each(func(_ int, opt *goquery.Selection) {
		val, exists := opt.Attr("value")
		if !exists || val == "0" || strings.TrimSpace(val) == "" {
			return
		}

		val = strings.Trim(val, "()")
		parts := strings.Split(val, ",")
		if len(parts) != 2 {
			return
		}

		lat, err1 := strconv.ParseFloat(strings.TrimSpace(parts[0]), 64)
		lon, err2 := strconv.ParseFloat(strings.TrimSpace(parts[1]), 64)
		if err1 != nil || err2 != nil {
			return
		}

		name := strings.TrimSpace(opt.Text())
		if name != "" {
			locations[name] = [2]float64{lat, lon}
		}
	})

	return locations
}

func parseStationRow(tr *goquery.Selection, locations map[string][2]float64) *BusStation {
	onclick, exists := tr.Attr("onclick")
	if !exists {
		return nil
	}

	// Extract "192" from "stop=192"
	matches := stopIDRegex.FindStringSubmatch(onclick)
	if len(matches) < 2 {
		return nil
	}
	code := matches[1] // this is "192"

	tds := tr.Find("td")
	if tds.Length() < 2 {
		return nil
	}

	// Extract displayed ID like "001"
	idStr := strings.TrimSpace(tds.Eq(1).Find("b").First().Text())
	name := strings.TrimSpace(tds.Eq(1).Find("b").Last().Text())

	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Fatalf("Error parsing ID %s: %s", idStr, err)
		return nil
	}

	lat, lon := 0.0, 0.0
	if coord, ok := locations[name]; ok {
		lat, lon = coord[0], coord[1]
	}

	// Store ID = 1, Code = "192"
	return NewBusStation(id, code, name, lat, lon)
}

func (p *HTMLParser) ParseBusStationDetails(html []byte) (*BusStationDetails, error) {
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(html))
	if err != nil {
		return nil, err
	}

	details := &BusStationDetails{}

	// Extract station ID
	doc.Find("#ModalBodyStopInfo table").First().Find("tr").Each(func(i int, s *goquery.Selection) {
		cells := s.Find("td")
		if cells.Length() != 2 {
			return
		}
		label := strings.TrimSpace(cells.Eq(0).Text())
		value := strings.TrimSpace(cells.Eq(1).Text())

		if strings.Contains(label, "Številka") {
			if id, err := strconv.Atoi(value); err == nil {
				details.ID = id
			}
		}
	})

	// Extract image URL
	if imgSrc, exists := doc.Find("#ModalBodyStopInfo img").Attr("src"); exists {
		details.ImageURL = strings.TrimSpace(imgSrc)
	}

	// Extract lines from the Linije row
	doc.Find("#ModalBodyStopInfo tr").Each(func(i int, s *goquery.Selection) {
		cells := s.Find("td")
		if cells.Length() != 2 {
			return
		}
		label := strings.TrimSpace(cells.Eq(0).Text())
		if label == "Linije" {
			cells.Eq(1).Find("a").Each(func(i int, a *goquery.Selection) {
				line := strings.TrimSpace(a.Text())
				if line != "" {
					details.Lines = append(details.Lines, line)
				}
			})
		}
	})

	// Extract departures from each <div id="l-XXX"> block
	doc.Find("div[id^='l-']").Each(func(i int, div *goquery.Selection) {
		lineID, exists := div.Attr("id")
		if !exists || !strings.HasPrefix(lineID, "l-") {
			return
		}
		line := strings.TrimPrefix(lineID, "l-")

		// Get the next <table> after this <div>
		table := div.NextAllFiltered("table").First()
		if table.Length() == 0 {
			log.Printf("No table found for line %s", line)
			return
		}

		rows := table.Find("tr")

		// Skip header (assumed first row)
		rows.Slice(1, goquery.ToEnd).Each(func(i int, row *goquery.Selection) {
			cells := row.Find("td")
			if cells.Length() != 2 {
				return
			}

			direction := strings.TrimSpace(cells.Eq(0).Text())
			rawTimes := strings.TrimSpace(cells.Eq(1).Text())

			if direction == "" || rawTimes == "" {
				return
			}

			directionLower := strings.ToLower(direction)
			if directionLower == "smer" || strings.Contains(strings.ToLower(rawTimes), "odhodi") {
				return // skip header row
			}

			times := strings.Fields(rawTimes)

			details.Departures = append(details.Departures, Departure{
				Line:      line,
				Direction: direction,
				Times:     times,
			})
		})
	})

	return details, nil
}
