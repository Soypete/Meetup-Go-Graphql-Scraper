package main

import (
	"fmt"
	"log"

	"github.com/soypete/Metup-Go-Graphql-Scraper/auth"
	"github.com/soypete/Metup-Go-Graphql-Scraper/meetup"
)

func main() {
	bearerToken, err := auth.GetBearerToken()
	if err != nil {
		log.Fatal(err)
	}
	meetupClient := meetup.Setup(bearerToken)

	// TODO(soypete): create csv with relevant data.
	// for list of data in README.md
	err = buildDataSet(meetupClient)
	if err != nil {
		log.Fatal(err)
	}
}

type analytics struct {
	Groups map[string]groupAnalytics // key in the group's meetup ID
}
type groupAnalytics struct {
	ID     string
	Name   string
	Events map[string]eventAnalytics // key is the group's meetup ID
}

type eventAnalytics struct {
	ID      string
	Title   string
	Date    string
	RSVPNum int
}

func getGroups(meetupClient meetup.Client) map[string]groupAnalytics {
	var cursor string
	hasNextPage := true
	groupMap := make(map[string]groupAnalytics)
	for hasNextPage {
		groups, err := meetupClient.GetListOfGroups(cursor)
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
func buildDataSet(meetupClient meetup.Client) error {
	var a analytics
	groupMap := getGroups(meetupClient)
	var cursor string
	hasNextPage := true
	eventMap := make(map[string]eventAnalytics)
	for hasNextPage {
		events := meetupClient.GetListOfEvents(cursor)
		for _, e := range events.Data.ProNetwork.EventsSearch.Edges {
			var ok bool
			if _, ok = eventMap[e.Node.ID]; !ok {
				eventMap[e.Node.ID] = eventAnalytics{
					ID:    e.Node.ID,
					Title: e.Node.Title,
					Date:  e.Node.DateTime,
				}
			}
			var g groupAnalytics
			if g, ok = groupMap[e.Node.Group.ID]; !ok || g.Events == nil {
				g.Events = make(map[string]eventAnalytics)
			}
			g.Events[e.Node.ID] = eventAnalytics{
				ID:    e.Node.ID,
				Title: e.Node.Title,
				Date:  e.Node.DateTime,
			}
			groupMap[e.Node.Group.ID] = g
		}
		cursor = events.Data.ProNetwork.EventsSearch.PageInfo.EndCursor
		hasNextPage = events.Data.ProNetwork.EventsSearch.PageInfo.HasNextPage
	}
	a.Groups = groupMap
	fmt.Println(a.Groups)
	return nil
}
