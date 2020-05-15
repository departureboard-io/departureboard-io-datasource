package main

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/departureboard-io/departureboard-io-datasource/pkg/departureboardio"
	"github.com/google/go-cmp/cmp"
	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/grafana-plugin-sdk-go/data"
)

func Test_translateTimeRangeToTimeWindowAndOffset(t *testing.T) {
	// Use unchanging times for tests.
	now := time.Date(2020, time.May, 11, 2, 3, 4, 5, time.UTC)
	anHourAgo := time.Date(2020, time.May, 11, 1, 3, 4, 5, time.UTC)
	anHourAway := time.Date(2020, time.May, 11, 3, 3, 4, 5, time.UTC)
	twoHoursAway := time.Date(2020, time.May, 11, 4, 3, 4, 5, time.UTC)

	type args struct {
		currentTime time.Time
		timeRange   backend.TimeRange
	}
	tests := []struct {
		name           string
		args           args
		wantTimeWindow int
		wantTimeOffset int
		wantErr        bool
	}{
		{
			"a time range that does not extend into the future should error",
			args{
				currentTime: now,
				timeRange: backend.TimeRange{
					From: anHourAgo,
					To:   now,
				},
			},
			0,
			0,
			true,
		},
		{
			"a time range that extends from now until an hour in the future",
			args{
				currentTime: now,
				timeRange: backend.TimeRange{
					From: now,
					To:   anHourAway,
				},
			},
			60,
			0,
			false,
		},
		{
			"a time range that extends from an hour in the future until two hours in the future",
			args{
				currentTime: now,
				timeRange: backend.TimeRange{
					From: anHourAway,
					To:   twoHoursAway,
				},
			},
			60,
			60,
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotTimeWindow, gotTimeOffset, err := translateTimeRangeToTimeWindowAndOffset(tt.args.currentTime, tt.args.timeRange)
			if (err != nil) != tt.wantErr {
				t.Errorf("translateTimeRangeToTimeWindowAndOffset() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotTimeWindow != tt.wantTimeWindow {
				t.Errorf("translateTimeRangeToTimeWindowAndOffset() gotTimeWindow = %v, want %v", gotTimeWindow, tt.wantTimeWindow)
			}
			if gotTimeOffset != tt.wantTimeOffset {
				t.Errorf("translateTimeRangeToTimeWindowAndOffset() gotTimeOffset = %v, want %v", gotTimeOffset, tt.wantTimeOffset)
			}
		})
	}
}

func TestDepartureBoardIODataSource_QueryData(t *testing.T) {
	crsCodes := []string{"PAD", "HAY", "NRW", "CBG"}
	type fields struct {
		DepartureBoardIOClient departureboardio.DepartureBoardIOClient
	}
	type args struct {
		ctx context.Context
		req *backend.QueryDataRequest
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *backend.QueryDataResponse
		wantErr bool
	}{
		{
			"no queries should return no responses",
			fields{departureboardio.NewFakeClient(crsCodes)},
			args{context.Background(), &backend.QueryDataRequest{}},
			&backend.QueryDataResponse{Responses: backend.Responses{}},
			false,
		},
		{
			"query for departures without service details",
			fields{departureboardio.NewFakeClient(crsCodes)},
			args{context.Background(), &backend.QueryDataRequest{
				Queries: []backend.DataQuery{{
					RefID:     "qfdwosd",
					QueryType: "",
					TimeRange: backend.TimeRange{
						From: time.Now(),
						To:   time.Now().Add(time.Hour),
					},
					JSON: func() []byte {
						b, err := json.Marshal(DepartureBoardIOQuery{
							StationCRS:     crsCodes[0],
							Arrivals:       false,
							Departures:     true,
							FilterCRS:      "",
							ServiceDetails: false,
						})
						if err != nil {
							panic(err)
						}
						return b
					}(),
				}},
			}},
			&backend.QueryDataResponse{Responses: backend.Responses{
				"qfdwosd": backend.DataResponse{
					Frames: data.Frames{
						data.NewFrame("PADDepartures",
							data.NewField("Scheduled", data.Labels{}, []string{"13:10"}),
							data.NewField("Estimated", data.Labels{}, []string{"On time"}),
							data.NewField("Destination", data.Labels{}, []string{"DAP"}),
							data.NewField("Platform", data.Labels{}, []string{"12A"}),
						)},
					Error: nil,
				},
			}},
			false,
		},
		{
			"query for departures with service details",
			fields{departureboardio.NewFakeClient(crsCodes)},
			args{context.Background(), &backend.QueryDataRequest{
				Queries: []backend.DataQuery{{
					RefID:     "qfdwsd",
					QueryType: "",
					TimeRange: backend.TimeRange{
						From: time.Now(),
						To:   time.Now().Add(time.Hour),
					},
					JSON: func() []byte {
						b, err := json.Marshal(DepartureBoardIOQuery{
							StationCRS:     crsCodes[0],
							Arrivals:       false,
							Departures:     true,
							FilterCRS:      "",
							ServiceDetails: true,
						})
						if err != nil {
							panic(err)
						}
						return b
					}(),
				}},
			}},
			&backend.QueryDataResponse{Responses: backend.Responses{
				"qfdwsd": backend.DataResponse{
					Frames: data.Frames{
						data.NewFrame("PADDepartures",
							data.NewField("Scheduled", data.Labels{}, []string{"13:10"}),
							data.NewField("Estimated", data.Labels{}, []string{"On time"}),
							data.NewField("Destination", data.Labels{}, []string{"DAP"}),
							data.NewField("Platform", data.Labels{}, []string{"12A"}),
							data.NewField("Service Details", data.Labels{}, []string{"PAD"}),
						)},
					Error: nil,
				},
			}},
			false,
		},
		{
			"query for arrivals without service details",
			fields{departureboardio.NewFakeClient(crsCodes)},
			args{context.Background(), &backend.QueryDataRequest{
				Queries: []backend.DataQuery{{
					RefID:     "qfawosd",
					QueryType: "",
					TimeRange: backend.TimeRange{
						From: time.Now(),
						To:   time.Now().Add(time.Hour),
					},
					JSON: func() []byte {
						b, err := json.Marshal(DepartureBoardIOQuery{
							StationCRS:     crsCodes[0],
							Arrivals:       true,
							Departures:     false,
							FilterCRS:      "",
							ServiceDetails: false,
						})
						if err != nil {
							panic(err)
						}
						return b
					}(),
				}},
			}},
			&backend.QueryDataResponse{Responses: backend.Responses{
				"qfawosd": backend.DataResponse{
					Frames: data.Frames{
						data.NewFrame("PADArrivals",
							data.NewField("Scheduled", data.Labels{}, []string{"13:10"}),
							data.NewField("Estimated", data.Labels{}, []string{"On time"}),
							data.NewField("Origin", data.Labels{}, []string{"DAP"}),
							data.NewField("Platform", data.Labels{}, []string{"12A"}),
						)},
					Error: nil,
				},
			}},
			false,
		},
		{
			"query for arrivals with service details",
			fields{departureboardio.NewFakeClient(crsCodes)},
			args{context.Background(), &backend.QueryDataRequest{
				Queries: []backend.DataQuery{{
					RefID:     "qfawsd",
					QueryType: "",
					TimeRange: backend.TimeRange{
						From: time.Now(),
						To:   time.Now().Add(time.Hour),
					},
					JSON: func() []byte {
						b, err := json.Marshal(DepartureBoardIOQuery{
							StationCRS:     crsCodes[0],
							Arrivals:       true,
							Departures:     false,
							FilterCRS:      "",
							ServiceDetails: true,
						})
						if err != nil {
							panic(err)
						}
						return b
					}(),
				}},
			}},
			&backend.QueryDataResponse{Responses: backend.Responses{
				"qfawsd": backend.DataResponse{
					Frames: data.Frames{
						data.NewFrame("PADArrivals",
							data.NewField("Scheduled", data.Labels{}, []string{"13:10"}),
							data.NewField("Estimated", data.Labels{}, []string{"On time"}),
							data.NewField("Origin", data.Labels{}, []string{"DAP"}),
							data.NewField("Platform", data.Labels{}, []string{"12A"}),
							data.NewField("Service Details", data.Labels{}, []string{"PAD"}),
						)},
					Error: nil,
				},
			}},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ds := &DepartureBoardIODataSource{
				DepartureBoardIOClient: tt.fields.DepartureBoardIOClient,
			}
			got, err := ds.QueryData(tt.args.ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("DepartureBoardIODataSource.QueryData() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if diff := cmp.Diff(tt.want, got, data.FrameTestCompareOptions()...); diff != "" {
				t.Errorf("DepartureBoardIODataSource.QueryData() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
