package meetup

import (
	"testing"
)

func TestGetAnalyticsFuncs(t *testing.T) {
	client := Setup("1234", "test", "csv")
	tests := []struct {
		name    string
		m       *Client
		funcs   string
		want    []func()
		wantErr bool
	}{
		{
			name:    "valid funcs",
			funcs:   "[groups, eventRSVP]",
			m:       &client,
			want:    []func(){client.GetGroupData, client.GetEventRSVPData},
			wantErr: false,
		},
		{
			name:    "invalid funcs",
			funcs:   "[groups, event-rsvp]",
			m:       &client,
			want:    nil,
			wantErr: true,
		},
		{
			name:    "no funcs provided",
			funcs:   "",
			m:       &client,
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.m.GetAnalyticsFunc(tt.funcs)
			if err != nil && !tt.wantErr {
				t.Errorf("GetAnalyticsFuncs() error = %v, wantErr %v", err, tt.wantErr)
			}
			if len(got) != len(tt.want) {
				t.Errorf("len of funcs = %d, want %d", len(got), len(tt.want))
			}
		})
	}
}
