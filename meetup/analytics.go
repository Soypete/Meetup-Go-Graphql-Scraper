package meetup

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type groupAnalytics struct {
	ID   string
	Name string
}

func (m Client) getGroups() map[string]groupAnalytics {
	var cursor string
	hasNextPage := true
	groupMap := make(map[string]groupAnalytics)
	for hasNextPage {
		groups, err := m.GetListOfGroups(cursor)
		if err != nil {
			log.Fatal(err)
		}

		for _, g := range groups.Data.ProNetwork.GroupsSearch.Edges {
			if _, ok := groupMap[g.Node.ID]; ok {
				continue
			}
			groupMap[g.Node.ID] = groupAnalytics{
				ID:   g.Node.ID,
				Name: g.Node.Name,
			}
		}
		cursor = groups.Data.ProNetwork.GroupsSearch.PageInfo.EndCursor
		hasNextPage = groups.Data.ProNetwork.GroupsSearch.PageInfo.HasNextPage
	}
	return groupMap
}

// TODO(soypete): maybe move the url here as a parameter to this function. The url is
// coupled to the auth id and token, so it would also make sense as an environment variable.

// BuildDataSet compiles the meetup.com groups and events data from a csv file used for analytics
func (m Client) BuildDataSet() error {
	groupMap := m.getGroups()
	// eventMap := make(map[string]eventAnalytics)
	csvFile, err := os.Create("temp-meetup.csv")
	if err != nil {
		return fmt.Errorf("failed to create csv, %w", err)
	}
	defer csvFile.Close()
	csvWriter := csv.NewWriter(csvFile)
	// write csv column names
	err = csvWriter.Write([]string{"event_id", "title", "date", "rsvp_num_going", "rsvp_num_waitlist", "group_id", "group_name"})
	if err != nil {
		return fmt.Errorf("could not write headers to csv file, %w", err)
	}
	defer csvWriter.Flush()

	var cursor string
	hasNextPage := true
	for hasNextPage {
		events := m.GetListOfEvents(cursor)
		for _, e := range events.Data.ProNetwork.EventsSearch.Edges {
			// write to the DuckDB
			formatWithSeconds := "2006-01-02T15:04:05-07:00"
			formatWithoutSeconds := "2006-01-02T15:04-07:00"
			// formatWithSecondsUTC := "2006-01-02T15:04:05Z"
			formatWithoutSecondsUTC := "2006-01-02T15:04Z"

			var eventDate time.Time
			var err error
			// compile to regexp if this works
			if strings.Contains(e.Node.DateTime, "Z") {
				eventDate, err = time.Parse(formatWithoutSecondsUTC, e.Node.DateTime)
				if err != nil {
					log.Printf("failed to parse date for event %s, %v\n", e.Node.ID, err)
				}
			} else {
				re := regexp.MustCompile(`\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}-\d{2}:\d{2}`) // 2019-02-07T18:00:00-05:00
				if re.MatchString(e.Node.DateTime) {
					eventDate, err = time.Parse(formatWithSeconds, e.Node.DateTime)
					if err != nil {
						log.Printf("failed to parse date for event %s, %v\n", e.Node.ID, err)
					}
				} else {
					eventDate, err = time.Parse(formatWithoutSeconds, e.Node.DateTime)
					if err != nil {
						log.Printf("failed to parse date for event %s, %v\n", e.Node.ID, err)
					}
				}
			}

			err = m.warehouse.InsertEvent(e.Node.ID, e.Node.Title, e.Node.Group.ID, e.Node.Group.Name, eventDate, e.Node.Going, e.Node.Waiting)
			if err != nil {
				log.Printf("failed to insert event %s, %v\n", e.Node.ID, err)
			}

			err = csvWriter.Write([]string{e.Node.ID, e.Node.Title, e.Node.DateTime, strconv.Itoa(e.Node.Going), strconv.Itoa(e.Node.Waiting), e.Node.Group.ID, groupMap[e.Node.Group.ID].Name})
			if err != nil {
				log.Printf("failed to write data for event %s, %v", e.Node.ID, err)
			}
		}
		cursor = events.Data.ProNetwork.EventsSearch.PageInfo.EndCursor
		hasNextPage = events.Data.ProNetwork.EventsSearch.PageInfo.HasNextPage
		csvWriter.Flush() // flushed after the pagination set
	}
	return nil
}
