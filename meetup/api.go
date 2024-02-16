package meetup

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

// tooling for interacting with meetup graphql api
type payloadql struct {
	Query     string `json:"query"`
	Variables string `json:"variables"`
}

type ProNetworkByUrlname struct {
	Data struct {
		ProNetwork struct {
			GroupsSearch GroupsSearch `json:"groupsSearch,omitempty"`
			EventsSearch EventsSearch `json:"eventsSearch,omitempty"`
		} `json:"proNetworkByUrlname"`
	} `json:"data"`
}
type GroupsSearch struct {
	Count    int `json:"count"`
	PageInfo struct {
		HasNextPage bool   `json:"hasNextPage"`
		StartCursor string `json:"startCursor"`
		EndCursor   string `json:"endCursor"`
	} `json:"pageInfo"`
	Edges []Edge `json:"edges"`
}

type Edge struct {
	Node Node `json:"node"`
}
type Node struct {
	ID    string `json:"id"`
	Title string `json:"title,omitempty"`
	Name  string `json:"name,omitempty"`
	Group struct {
		ID   string `json:"id,omitempty"`
		Name string `json:"name,omitempty"`
	} `json:"group,omitempty"`
	DateTime string `json:"dateTime,omitempty"`
	Going    int    `json:"going,omitempty"`
	Waiting  int    `json:"waiting,omitempty"`
}

type EventsSearch struct {
	Count    int `json:"count"`
	PageInfo struct {
		HasNextPage bool   `json:"hasNextPage"`
		StartCursor string `json:"startCursor"`
		EndCursor   string `json:"endCursor"`
	} `json:"pageInfo"`
	Edges []Edge `json:"edges"`
}

// TODO(soypete): edit variables
func getInputandVariables(isFirst bool, lastCursor, urlname string, numPerPage int) (string, string) {
	if isFirst {
		return "input: {first: $itemsNum}", fmt.Sprintf(`{"urlname":"%s", "itemsNum":%d}`, urlname, numPerPage)
	}
	return "input: {first: $itemsNum, after: $cursor}", fmt.Sprintf(`{"urlname":"%s", "itemsNum":%d, "cursor":"%s"}`, urlname, numPerPage, lastCursor)

}

var queryTemplate = `query (%s) { proNetworkByUrlname(urlname: $urlname) { %s(%s) {count pageInfo { hasNextPage startCursor endCursor } edges { node { id %s } } } }}`

func makePayloadql(isGroup, isfirst bool, lastCursor, urlname string, numPerPage int) payloadql {
	variableTypes := "$urlname: String!, $itemsNum: Int!"
	if !isfirst {
		variableTypes = variableTypes + ", $cursor: String!"
	}
	searchType := `eventsSearch`
	nodeQuery := `title group { id name } dateTime going waiting`
	if isGroup {
		searchType = `groupsSearch`
		nodeQuery = `name`
	}

	input, variables := getInputandVariables(isfirst, lastCursor, urlname, numPerPage)
	query := fmt.Sprintf(queryTemplate, variableTypes, searchType, input, nodeQuery)
	p := payloadql{
		Query:     query,
		Variables: variables,
	}
	return p
}
func (c Client) getListOfGroups(cursor string) (ProNetworkByUrlname, error) {
	isFirst := true
	if cursor != "" {
		isFirst = false
	}
	p := makePayloadql(true, isFirst, cursor, c.proname, 3)
	body, err := c.sendRequest(p)
	if err != nil {
		return ProNetworkByUrlname{}, err
	}
	var respData ProNetworkByUrlname
	err = json.Unmarshal(body, &respData)
	if err != nil {
		return ProNetworkByUrlname{}, err
	}
	return respData, nil
}

type Group struct {
	Data struct {
		Group struct {
			ID             string `json:"id"`
			Name           string `json:"name"`
			URLKey         string `json:"urlname"`
			Link           string `json:"link"`
			Status         string `json:"status"`
			GroupAnalytics struct {
				totalMembers         int     `json:"totalMembers,omitempty"`
				lastEventDate        string  `json:"lastEventDate,omitempty"`
				averageAge           float64 `json:"averageAge,omitempty"`
				totalPastEvents      int     `json:"totalPastEvents,omitempty"`
				totalPastRSVPs       int     `json:"totalPastRSVPs,omitempty"`
				repeatRSVPers        int     `json:"repeatRSVPers,omitempty"`
				averageRSVPsPerEvent float64 `json:"averageRSVPsPerEvent,omitempty"`
				totalUpcomingEvents  int     `json:"totalUpcomingEvents,omitempty"`
			} `json:"groupAnalytics,omitempty"`
		} `json:"group"`
	} `json:"data"`
}

func (c Client) getGroupByID(groupID string) Group {
	var groupQueryTemplate = `query ($groupID: ID) {
  group(id: $groupID) {
    id
    name
    urlname
    link
    status
    groupAnalytics {
      totalMembers
      lastEventDate
      averageAge
      totalPastEvents
      totalPastRsvps
      totalRepeatRsvpers
      averageRsvpsPerEvent
      totalUpcomingEvents
    }
  }
} `
	p := payloadql{
		Query:     groupQueryTemplate,
		Variables: fmt.Sprintf(`{"groupID":"%s"}`, groupID),
	}
	body, err := c.sendRequest(p)
	if err != nil {
		log.Fatal(err)
	}
	var respData Group
	// run it and capture the response
	err = json.Unmarshal(body, &respData)
	if err != nil {
		log.Fatal(err)
	}
	return respData
}

func (c Client) getListOfEvents(cursor string) ProNetworkByUrlname {
	isFirst := true
	if cursor != "" {
		isFirst = false
	}
	p := makePayloadql(false, isFirst, cursor, c.proname, 3)
	body, err := c.sendRequest(p)
	if err != nil {
		log.Fatal(err)
	}
	// run it and capture the response
	var respData ProNetworkByUrlname
	err = json.Unmarshal(body, &respData)
	if err != nil {
		log.Fatal(err)
	}
	return respData
}

func (c Client) sendRequest(ql payloadql) (resp []byte, err error) {
	b, err := json.Marshal(ql)
	if err != nil {
		return nil, fmt.Errorf("cannot marshal graphql payload, %w", err)
	}
	req, err := http.NewRequest("POST", c.url, bytes.NewReader(b))
	if err != nil {
		return nil, fmt.Errorf("failed to create new http.Request, %w", err)
	}
	// set header fields
	req.Header.Add("Authorization", c.bearerToken)
	req.Header.Add("Content-Type", "application/json")

	res, err := c.ql.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed meetup.com api request, %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode >= 300 {
		return nil, fmt.Errorf("meetup.com api request failed with status: %s", res.Status)
	}

	// TODO(soypete): check payload status
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read payload body, %w", err)
	}
	return body, nil
}
