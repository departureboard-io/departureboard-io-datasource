//+build e2e

package departureboardio

import (
	"net/http"
	"net/url"
	"os"
	"testing"
)

// TestDepartureBoardIO_GetDeparturesByCRS currently only tests to see that we got a board-like response from the API.
// It mostly just validates that the right URL is requested.
func TestDepartureBoardIO_GetDeparturesByCRS(t *testing.T) {
	type fields struct {
		Client      http.Client
		APIEndpoint url.URL
		APIKey      string
	}
	type args struct {
		crs     string
		options boardOptions
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *Board
		wantErr bool
	}{
		{
			"PAD",
			fields{
				Client: http.Client{},
				APIEndpoint: func() url.URL {
					parsed, _ := url.Parse("https://api.departureboard.io/api/v2.0")
					return *parsed
				}(),
				APIKey: os.Getenv("NATIONALRAIL_API_KEY"),
			},
			args{
				crs:     "PAD",
				options: NewDefaultBoardOptions(),
			},
			&Board{},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Client{
				Client:      tt.fields.Client,
				APIEndpoint: tt.fields.APIEndpoint,
				APIKey:      tt.fields.APIKey,
			}
			got, err := c.GetDeparturesByCRS(tt.args.crs, tt.args.options)
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
