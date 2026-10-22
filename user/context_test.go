package user

import (
	"context"
	"testing"
)

func TestContextWithID(t *testing.T) {
	tests := []struct {
		name    string
		id      string
		wantErr bool
	}{
		{name: "valid", id: "u_123", wantErr: false},
		{name: "empty", id: "", wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, err := ContextWithID(context.Background(), tt.id)
			if (err != nil) != tt.wantErr {
				t.Fatalf("ContextWithID(%q) error = %v, wantErr %v", tt.id, err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}

			got, ok := IDFromContext(ctx)
			if !ok {
				t.Fatal("IDFromContext() ok = false, want true")
			}
			if got != tt.id {
				t.Errorf("IDFromContext() = %q, want %q", got, tt.id)
			}
		})
	}
}

func TestIDFromContext_Absent(t *testing.T) {
	if _, ok := IDFromContext(context.Background()); ok {
		t.Fatal("IDFromContext() on empty context: ok = true, want false")
	}
}
