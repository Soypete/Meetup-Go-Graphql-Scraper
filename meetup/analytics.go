package meetup

import (
	"encoding/csv"
	"log"
	"os"
	"strconv"
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
		groups, err := m.getListOfGroups(cursor)
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

func (m Client) GetGroupData() {
	groupMap := m.getGroups()
	csvFile, err := os.Create("meetup_groups.csv")
	if err != nil {
		log.Printf("failed to create csv, %v", err)
	}
	csvWriter := csv.NewWriter(csvFile)
	err = csvWriter.Write([]string{"group_id", "group_name"})
	if err != nil {
		log.Printf("could not write headers to csv file, %v", err)
		return
	}
	defer csvWriter.Flush()
	for id, g := range groupMap {
		err := csvWriter.Write([]string{id, g.Name})
		if err != nil {
			log.Printf("failed to write data for group %s, %v", id, err)
		}
	}
}

// GetEventRSVPData compiles the meetup.com groups and events data from a csv file used for analytics
func (m Client) GetEventRSVPData() {
	groupMap := m.getGroups()
	csvFile, err := os.Create("event_RSVP.csv")
	if err != nil {
		log.Printf("failed to create csv, %v", err)
		return
	}
	defer csvFile.Close()
	csvWriter := csv.NewWriter(csvFile)
	// write csv column names
	err = csvWriter.Write([]string{"event_id", "title", "date", "rsvp_num_going", "rsvp_num_waitlist", "group_id", "group_name"})
	if err != nil {
		log.Printf("could not write headers to csv file, %v", err)
		return
	}

	var cursor string
	hasNextPage := true
	for hasNextPage {
		events := m.getListOfEvents(cursor)
		for _, e := range events.Data.ProNetwork.EventsSearch.Edges {
			err := csvWriter.Write([]string{e.Node.ID, e.Node.Title, e.Node.DateTime, strconv.Itoa(e.Node.Going), strconv.Itoa(e.Node.Waiting), e.Node.Group.ID, groupMap[e.Node.Group.ID].Name})
			if err != nil {
				log.Printf("failed to write data for event %s, %v", e.Node.ID, err)
			}
		}
		cursor = events.Data.ProNetwork.EventsSearch.PageInfo.EndCursor
		hasNextPage = events.Data.ProNetwork.EventsSearch.PageInfo.HasNextPage
		csvWriter.Flush() // flushed after the pagination set
	}
}
