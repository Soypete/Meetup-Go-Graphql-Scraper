// meetup module contains all the logic for making our
// data requests to Meetup.com's grapql-api.
// The api docs can be found here: https://www.meetup.com/api/guide/#graphQl-guide.
package meetup

import (
	"fmt"
	"net/http"
	"strings"
	"time"
)

// Client connects all information for connecting to the Meetup API.
type Client struct {
	proname     string
	ql          *http.Client
	bearerToken string //TODO(soypete): refresh
	url         string
}

var (
	analyticsFuncs map[string]func()
)

func Setup(bearerToken string, proname string) Client {
	// TODO: add check for valid bearerToken
	analyticsFuncs = make(map[string]func())
	c := Client{
		proname: proname,
		ql: &http.Client{
			Timeout: 10 * time.Second,
		},
		bearerToken: bearerToken,
		url:         "https://api.meetup.com/gql",
	}
	analyticsFuncs["groups"] = c.GetGroupData
	analyticsFuncs["eventRSVP"] = c.GetEventRSVPData
	return c
}

func (m Client) GetAnalyticsFunc(funcs string) ([]func(), error) {
	funcs = strings.Trim(funcs, "[]")
	funcsList := strings.Split(funcs, ",")
	var queuedFuncs []func()
	for _, f := range funcsList {
		key := strings.Trim(f, " ")
		if _, ok := analyticsFuncs[key]; !ok {
			return nil, fmt.Errorf("invalid flag: %s", key)
			// TODO: return an error for invalid flad
		}
		queuedFuncs = append(queuedFuncs, analyticsFuncs[key])
	}
	return queuedFuncs, nil
}
