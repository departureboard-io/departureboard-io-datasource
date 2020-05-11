package main

import (
	"testing"
	"time"

	"github.com/grafana/grafana-plugin-sdk-go/backend"
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
