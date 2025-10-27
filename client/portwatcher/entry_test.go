package portwatcher

import (
	"reflect"
	"testing"
)

func TestSnapshotDiff(t *testing.T) {
	prev := Snapshot{
		Entries: []Entry{
			{Port: 8080, Addr: "127.0.0.1", Protocol: "tcp", Process: "app"},
			{Port: 9090, Addr: "0.0.0.0", Protocol: "tcp", Process: "svc"},
		},
	}
	curr := Snapshot{
		Entries: []Entry{
			{Port: 8080, Addr: "127.0.0.1", Protocol: "tcp", Process: "app"},
			{Port: 3030, Addr: "127.0.0.1", Protocol: "tcp", Process: "new"},
			{Port: 3030, Addr: "127.0.0.1", Protocol: "tcp", Process: "new"}, // duplicate
		},
	}

	diff := curr.Diff(prev)
	wantAdded := []Entry{
		{Port: 3030, Addr: "127.0.0.1", Protocol: "tcp", Process: "new"},
	}
	wantRemoved := []Entry{
		{Port: 9090, Addr: "0.0.0.0", Protocol: "tcp", Process: "svc"},
	}

	if !reflect.DeepEqual(diff.Added, wantAdded) {
		t.Fatalf("diff.Added = %#v, want %#v", diff.Added, wantAdded)
	}
	if !reflect.DeepEqual(diff.Removed, wantRemoved) {
		t.Fatalf("diff.Removed = %#v, want %#v", diff.Removed, wantRemoved)
	}
}

func TestConfigShouldInclude(t *testing.T) {
	cfg := Config{
		MinPort:          1000,
		MaxPort:          9000,
		ExcludePorts:     []uint16{3306},
		ExcludeProcesses: []string{"ignore-me"},
	}

	tests := map[string]struct {
		entry Entry
		want  bool
	}{
		"below min port": {
			entry: Entry{Port: 80, Process: "foo"},
			want:  false,
		},
		"above max port": {
			entry: Entry{Port: 10001, Process: "foo"},
			want:  false,
		},
		"excluded port": {
			entry: Entry{Port: 3306, Process: "foo"},
			want:  false,
		},
		"excluded process": {
			entry: Entry{Port: 5000, Process: "IGNORE-ME"},
			want:  false,
		},
		"allowed entry": {
			entry: Entry{Port: 5000, Process: "service"},
			want:  true,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			if got := cfg.ShouldInclude(tt.entry); got != tt.want {
				t.Fatalf("ShouldInclude(%#v) = %v, want %v", tt.entry, got, tt.want)
			}
		})
	}
}
