package main

import (
	"testing"
	"time"
)

func TestHumanData(t *testing.T) {

	tests := []struct {
		name string
		tm   time.Time
		want string
	}{
		{
			name: "standard time",
			tm:   time.Date(2026, 5, 27, 17, 39, 33, 19, time.UTC),
			want: "27 May 2026 at 17:39",
		},
		{
			name: "zero time",
			tm:   time.Time{},
			want: "",
		},
		{
			name: "time zone: CET (UTC+1)",
			tm:   time.Date(2026, 5, 27, 17, 39, 33, 19, time.FixedZone("CET", 1*60*60)),
			want: "27 May 2026 at 16:39",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hd := humanDate(tt.tm)

			if hd != tt.want {
				t.Errorf("got %q, want %q", hd, tt.want)
			}
		})
	}

	// tm := time.Date(2026, 5, 27, 17, 39, 33, 19, time.UTC)
	// expect := "27 May 2026 at 17:39"
	// hd := humanDate(tm)

	// if hd != expect {
	// 	t.Errorf("got %q, want %q", hd, expect)
	// }
}
