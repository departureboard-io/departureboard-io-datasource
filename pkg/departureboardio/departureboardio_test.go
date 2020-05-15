//+build e2e

package departureboardio

import (
	"net/http"
	"os"
	"testing"
)

// TestDepartureBoardIO_GetDeparturesByCRS currently only tests to see that we got a board-like response from the API.
// It mostly just validates that the right URL is requested.
func TestDepartureBoardIO_GetDeparturesByCRS(t *testing.T) {
	type fields struct {
		Client http.Client
	}
	type args struct {
		apiEndpoint string
		apiKey      string
		crs         string
		options     boardOptions
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *DepartureBoard
		wantErr bool
	}{
		{
			"PAD",
			fields{
				Client: http.Client{},
			},
			args{
				apiEndpoint: "https://api.departureboard.io/api/v2.0",
				apiKey:      os.Getenv("NATIONALRAIL_API_KEY"),
				crs:         "PAD",
				options:     NewDefaultBoardOptions(),
			},
			&DepartureBoard{},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Client{
				Client: tt.fields.Client,
			}
			got, err := c.GetDeparturesByCRS(tt.args.apiEndpoint, tt.args.apiKey, tt.args.crs, tt.args.options)
			if (err != nil) != tt.wantErr {
				t.Errorf("DepartureBoardIO.GetDeparturesByCRS() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got == nil {
				t.Errorf("DepartureBoardIO.GetDeparturesByCRS() = %v, want %v", got, tt.want)
			}
		})
	}
}
