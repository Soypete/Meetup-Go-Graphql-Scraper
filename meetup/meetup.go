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
	bearerToken string
	filesDir    string
	url         string
}

var (
	analyticsFuncs map[string]func()
)

// Setup initializes the Meetup client with impletented analytics functions.
// The bearerToken is the token used to authenticate with the Meetup API.
// The proname is the name of the project that is using the Meetup API.
// the client is returned with all credentials and implemented analytics functions needed to
// create analytics files.
func Setup(bearerToken, proname, directory string) Client {
	// TODO: add check for valid bearerToken
	analyticsFuncs = make(map[string]func())
	c := Client{
		proname: proname,
		ql: &http.Client{
			Timeout: 10 * time.Second,
		},
		filesDir:    directory,
		bearerToken: bearerToken,
		url:         "https://api.meetup.com/gql",
	}
	analyticsFuncs["groups"] = c.GetGroupList
	analyticsFuncs["groupAnalytics"] = c.GetGroupData
	analyticsFuncs["eventRSVP"] = c.GetEventRSVPData
	return c
}

// GetAnalyticsFunc returns a list of analytics functions that provided
// via the command line flag.
func (m Client) GetAnalyticsFunc(funcs string) ([]func(), error) {
	funcs = strings.Trim(funcs, "[]")
	funcsList := strings.Split(funcs, ",")
	var queuedFuncs []func()
	for _, f := range funcsList {
		key := strings.Trim(f, " ")
		if _, ok := analyticsFuncs[key]; !ok {
			return nil, fmt.Errorf("invalid flag: %s", key)
		}
		queuedFuncs = append(queuedFuncs, analyticsFuncs[key])
	}
	return queuedFuncs, nil
}
