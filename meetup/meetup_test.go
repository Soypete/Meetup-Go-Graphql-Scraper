// meetup module contains all the logic for making our
// data requests to Meetup.com's grapql-api.
// The api docs can be found here: https://www.meetup.com/api/guide/#graphQl-guide.
package meetup

import (
	"testing"
)

func Test_getInputandVariables(t *testing.T) {
	type args struct {
		isFirst    bool
		lastCursor string
		urlname    string
		numPerPage int
	}
	tests := []struct {
		name          string
		args          args
		wantInput     string
		wantVariables string
	}{
		// TODO: Add test cases.
		{name: "is first",
			args: args{
				isFirst:    true,
				urlname:    "go",
				numPerPage: 3,
			},
			wantInput:     "input: {first: $itemsNum}",
			wantVariables: `{"urlname":"go","itemsNum": 3}`,
		},
		{name: "has cursor",
			args: args{
				lastCursor: "1234",
				urlname:    "go",
				numPerPage: 3,
			},
			wantInput:     "input: {first: $itemsNum, after: $cursor}",
			wantVariables: `{"urlname":"go","itemsNum": 3,"cursor": "1234"}`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotInput, gotVariables := getInputandVariables(tt.args.isFirst, tt.args.lastCursor, tt.args.urlname, tt.args.numPerPage)
			if gotInput != tt.wantInput {
				t.Errorf("getInputandVariables() gotInput = %v, wantInput %v", gotInput, tt.wantInput)
			}
			if gotVariables != tt.wantVariables {
				t.Errorf("getInputandVariables() gotVariables = %v, want %v", gotVariables, tt.wantVariables)
			}
		})
	}
}

func Test_makePayloadql(t *testing.T) {
	type args struct {
		isGroup    bool
		isFirst    bool
		lastCursor string
		urlname    string
		numPerPage int
	}
	tests := []struct {
		name string
		args args
		want payloadql
	}{
		{name: "groups first query",
			args: args{
				isGroup:    true,
				isFirst:    true,
				urlname:    "go",
				numPerPage: 3,
			},
			want: payloadql{
				Variables: `{"urlname": "go", "itemsNum": 3}`,
				Query: `query ($urlname: String!, $itemsNum: Int!) {
  proNetworkByUrlname(urlname: $urlname) {
    groupsSearch(input: {first: $itemsNum}) {
      count
      pageInfo {
        hasNextPage
        startCursor
        endCursor
      }
      edges {
        node {
          id
          name
        }
      }
    }
  }
}`,
			},
		},
		{name: "events first query",
			args: args{
				isGroup:    true,
				isFirst:    true,
				urlname:    "go",
				numPerPage: 3,
			},
			want: payloadql{
				Variables: `{"urlname": "go", "itemsNum": 3}`,
				Query: `query ($urlname: String!) { 
		proNetworkByUrlname(urlname: $urlname) { 
			eventsSearch(input: {first: $itemsNum}) {
      count
      pageInfo {
				hasNextPage
				startCursor
        endCursor
      }
      edges {
        node {
          id
         	title
					group {
						id
						name
					}
					dateTime
				} 
			} 
		} 
	} 
}`,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := makePayloadql(tt.args.isGroup, tt.args.isFirst, tt.args.lastCursor, tt.args.urlname, tt.args.numPerPage); got.Variables == tt.want.Variables {
				t.Errorf("makePayloadql() = %v, want %v", got, tt.want)
			}
		})
	}
}
