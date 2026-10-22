package timezone

import (
	"context"
	"testing"

	// tzdata is embedded so IANA zone lookups are hermetic: these tests pass
	// on minimal images (scratch/distroless CI) that ship no system zoneinfo.
	_ "time/tzdata"
)

func TestContextWithZone(t *testing.T) {
	tests := []struct {
		name    string
		zone    string
		wantErr bool
	}{
		{name: "iana zone", zone: "Asia/Jakarta", wantErr: false},
		{name: "utc", zone: "UTC", wantErr: false},
		{name: "empty", zone: "", wantErr: true},
		{name: "local", zone: "Local", wantErr: true},
		{name: "invalid", zone: "Mars/Olympus", wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, err := ContextWithZone(context.Background(), tt.zone)
			if (err != nil) != tt.wantErr {
				t.Fatalf("ContextWithZone(%q) error = %v, wantErr %v", tt.zone, err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}

			got, ok := ZoneFromContext(ctx)
			if !ok {
				t.Fatal("ZoneFromContext() ok = false, want true")
			}
			if got.String() != tt.zone {
				t.Errorf("ZoneFromContext() = %q, want %q", got, tt.zone)
			}
		})
	}
}

func TestZoneFromContext_Absent(t *testing.T) {
	if _, ok := ZoneFromContext(context.Background()); ok {
		t.Fatal("ZoneFromContext() on empty context: ok = true, want false")
	}
}
