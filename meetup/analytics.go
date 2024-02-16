package meetup

import (
	"encoding/csv"
	"fmt"
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
	csvFile, err := os.Create("meetup_groups_analytics.csv")
	if err != nil {
		log.Printf("failed to create csv, %v", err)
	}

	defer csvFile.Close()
	csvWriter := csv.NewWriter(csvFile)
	err = csvWriter.Write([]string{"group_id", "group_name", "group_url_key", "url", "status", "totalMembers", "lastEventDate", "averageAge",
		"totalPastEvents", "totalPastRSVPs", "repeatRSVPers", "averageRSVPsPerEvent", "totalUpcomingEvents"})
	if err != nil {
		log.Printf("could not write headers to csv file, %v", err)
		return
	}
	defer csvWriter.Flush()
	for id, g := range groupMap {
		group := m.getGroupByID(id)
		err := csvWriter.Write([]string{id, g.Name, group.Data.Group.URLKey, group.Data.Group.Link, group.Data.Group.Status, strconv.Itoa(group.Data.Group.GroupAnalytics.totalMembers),
			group.Data.Group.GroupAnalytics.lastEventDate, fmt.Sprintf("%f", group.Data.Group.GroupAnalytics.averageAge), strconv.Itoa(group.Data.Group.GroupAnalytics.totalPastEvents),
			strconv.Itoa(group.Data.Group.GroupAnalytics.totalPastRSVPs), strconv.Itoa(group.Data.Group.GroupAnalytics.repeatRSVPers),
			fmt.Sprintf("%f", group.Data.Group.GroupAnalytics.averageRSVPsPerEvent), strconv.Itoa(group.Data.Group.GroupAnalytics.totalUpcomingEvents)})
		if err != nil {
			log.Printf("failed to write data for group %s, %v", id, err)
		}
	}
}

// GetGroupList returns a list of meetup.com groups in a csv file
func (m Client) GetGroupList() {
	groupMap := m.getGroups()
	csvFile, err := os.Create("meetup_groups.csv")
	if err != nil {
		log.Printf("failed to create csv, %v", err)
	}
	defer csvFile.Close()
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
