package meetup

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
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
		{
			name: "Test get Input and Variables",
			args: args{
				isFirst:    true,
				urlname:    "test",
				numPerPage: 10,
			},
			wantInput:     "input: {first: $itemsNum}",
			wantVariables: `{"urlname":"test", "itemsNum":10}`,
		},
		{
			name: "Test get Input and Variables with cursor",
			args: args{
				isFirst:    false,
				lastCursor: "timestamp",
				urlname:    "test",
				numPerPage: 10,
			},
			wantInput:     "input: {first: $itemsNum, after: $cursor}",
			wantVariables: `{"urlname":"test", "itemsNum":10, "cursor":"timestamp"}`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			input, variables := getInputandVariables(tt.args.isFirst, tt.args.lastCursor, tt.args.urlname, tt.args.numPerPage)
			if input != tt.wantInput {
				t.Errorf("getInputandVariables() \ngot = %v\n, want= %v\n", input, tt.wantInput)
			}
			if variables != tt.wantVariables {
				t.Errorf("getInputandVariables() got1 = %v, want %v", variables, tt.wantVariables)
			}
		})
	}
}

func Test_makePayloadql(t *testing.T) {
	type args struct {
		isGroup    bool
		isfirst    bool
		lastCursor string
		urlname    string
		numPerPage int
	}
	tests := []struct {
		name string
		args args
		want payloadql
	}{
		{
			name: "Test makeGroup Payload",
			args: args{
				isGroup:    true,
				isfirst:    true,
				urlname:    "test",
				numPerPage: 10,
			},
			want: payloadql{
				Query:     "query ($urlname: String!, $itemsNum: Int!) { proNetworkByUrlname(urlname: $urlname) { groupsSearch(input: {first: $itemsNum}) {count pageInfo { hasNextPage startCursor endCursor } edges { node { id name } } } }}",
				Variables: `{"urlname":"test", "itemsNum":10}`},
		},
		{
			name: "Test makeEvent Payload",
			args: args{
				isGroup:    false,
				isfirst:    true,
				urlname:    "test",
				numPerPage: 10,
			},
			want: payloadql{
				Query:     "query ($urlname: String!, $itemsNum: Int!) { proNetworkByUrlname(urlname: $urlname) { eventsSearch(input: {first: $itemsNum}) {count pageInfo { hasNextPage startCursor endCursor } edges { node { id title group { id name } dateTime going waiting } } } }}",
				Variables: `{"urlname":"test", "itemsNum":10}`},
		},
		{
			name: "Test make Group Payload with cursor",
			args: args{
				isGroup:    true,
				isfirst:    false,
				lastCursor: "timestamp",
				urlname:    "test",
				numPerPage: 10,
			},
			want: payloadql{
				Query:     "query ($urlname: String!, $itemsNum: Int!, $cursor: String!) { proNetworkByUrlname(urlname: $urlname) { groupsSearch(input: {first: $itemsNum, after: $cursor}) {count pageInfo { hasNextPage startCursor endCursor } edges { node { id name } } } }}",
				Variables: `{"urlname":"test", "itemsNum":10, "cursor":"timestamp"}`},
		},
		{
			name: "Test make Group Payload with Last Cursor",
			args: args{
				isGroup:    true,
				isfirst:    false,
				lastCursor: "timestamp",
				urlname:    "test",
				numPerPage: 10,
			},
			want: payloadql{
				Query:     "query ($urlname: String!, $itemsNum: Int!, $cursor: String!) { proNetworkByUrlname(urlname: $urlname) { groupsSearch(input: {first: $itemsNum, after: $cursor}) {count pageInfo { hasNextPage startCursor endCursor } edges { node { id name } } } }}",
				Variables: `{"urlname":"test", "itemsNum":10, "cursor":"timestamp"}`},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := makePayloadql(tt.args.isGroup, tt.args.isfirst, tt.args.lastCursor, tt.args.urlname, tt.args.numPerPage); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("makePayloadql() \ngot = %v\nwant= %v\n", got, tt.want)
			}
		})
	}
}

func TestClient_sendRequest(t *testing.T) {
	t.Run("Test sendRequest", testSuccess_sendRequest)
	t.Run("Test sendRequest with 400", testWithError_sendRequest)
}

func testSuccess_sendRequest(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "Hello, client")
	}))
	c := Client{
		ql:  ts.Client(),
		url: ts.URL,
	}
	ql := payloadql{
		Query:     "query ($urlname: String!, $itemsNum: Int!) { proNetworkByUrlname(urlname: $urlname) { eventsSearch(input: {first: $itemsNum}) {count pageInfo { hasNextPage startCursor endCursor } edges { node { id title group { id name } dateTime going waiting } } } }}",
		Variables: `{"urlname":"test", "itemsNum":10}`}

	gotResp, err := c.sendRequest(ql)
	if err != nil {
		t.Errorf("Client.sendRequest() error = %v", err)
		return
	}
	resp := string(gotResp)
	wantResp := "Hello, client"
	if !reflect.DeepEqual(resp, wantResp) {
		t.Errorf("Client.sendRequest() = %v, want %v", resp, wantResp)
	}
}
func testWithError_sendRequest(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "Sorry, client")
	}))
	c := Client{
		ql:  ts.Client(),
		url: ts.URL,
	}
	ql := payloadql{
		Query:     "query ($urlname: String!, $itemsNum: Int!) { proNetworkByUrlname(urlname: $urlname) { eventsSearch(input: {first: $itemsNum}) {count pageInfo { hasNextPage startCursor endCursor } edges { node { id title group { id name } dateTime going waiting } } } }}",
		Variables: `{"urlname":"test", "itemsNum":10}`}

	gotResp, err := c.sendRequest(ql)
	if err == nil {
		t.Errorf("Expecting Error on Client.sendRequest() but got none")
		return
	}
	if gotResp != nil {
		t.Errorf("Client.sendRequest() = %v, want nil", gotResp)
	}
}
