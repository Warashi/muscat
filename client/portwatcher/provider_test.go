package portwatcher

import (
	"context"
	"io"
	"log"
	"os"
	"testing"

	gnet "github.com/shirou/gopsutil/v4/net"
)

func TestSystemSnapshotProvider(t *testing.T) {
	ctx := context.Background()
	fetchCalls := 0
	resolverCalls := 0
	provider := &systemSnapshotProvider{
		cfg: Config{MinPort: 1000},
		fetcher: func(context.Context) ([]gnet.ConnectionStat, error) {
			fetchCalls++
			return []gnet.ConnectionStat{
				{
					Status: "LISTEN",
					Laddr:  gnet.Addr{IP: "127.0.0.1", Port: 80},
					Pid:    1111,
				},
				{
					Status: "LISTEN",
					Laddr:  gnet.Addr{IP: "127.0.0.1", Port: 3000},
					Pid:    2222,
				},
				{
					Status: "LISTEN",
					Laddr:  gnet.Addr{IP: "", Port: 4000},
					Pid:    2222,
				},
				{
					Status: "LISTEN",
					Laddr:  gnet.Addr{IP: "[::]", Port: 5000},
					Pid:    3333,
				},
				{
					Status: "ESTABLISHED",
					Laddr:  gnet.Addr{IP: "127.0.0.1", Port: 6000},
					Pid:    4444,
				},
			}, nil
		},
		processResolver: func(ctx context.Context, pid int32) (string, error) {
			resolverCalls++
			switch pid {
			case 2222:
				return "app", nil
			case 3333:
				return "svc", nil
			default:
				return "", nil
			}
		},
		logger: log.New(io.Discard, "", 0),
	}

	snapshot, err := provider.Snapshot(ctx)
	if err != nil {
		t.Fatalf("Snapshot returned error: %v", err)
	}
	if fetchCalls != 1 {
		t.Fatalf("expected fetcher to be called once, got %d", fetchCalls)
	}
	if resolverCalls != 2 {
		t.Fatalf("expected resolver to be called twice, got %d", resolverCalls)
	}

	wantEntries := []Entry{
		{Addr: "127.0.0.1", Port: 3000, Protocol: "tcp", Process: "app"},
		{Addr: "0.0.0.0", Port: 4000, Protocol: "tcp", Process: "app"},
		{Addr: "0.0.0.0", Port: 5000, Protocol: "tcp", Process: "svc"},
	}

	if len(snapshot.Entries) != len(wantEntries) {
		t.Fatalf("snapshot entries length = %d, want %d", len(snapshot.Entries), len(wantEntries))
	}
	for _, want := range wantEntries {
		found := false
		for _, entry := range snapshot.Entries {
			if entry == want {
				found = true
				break
			}
		}
		if !found {
			t.Fatalf("expected entry %#v not found in snapshot %#v", want, snapshot.Entries)
		}
	}
}

func TestSystemSnapshotProviderPermissionError(t *testing.T) {
	provider := &systemSnapshotProvider{
		cfg: Config{},
		fetcher: func(context.Context) ([]gnet.ConnectionStat, error) {
			return nil, os.ErrPermission
		},
		processResolver: func(context.Context, int32) (string, error) {
			return "", nil
		},
		logger: log.New(io.Discard, "", 0),
	}

	snapshot, err := provider.Snapshot(context.Background())
	if err != nil {
		t.Fatalf("Snapshot returned error: %v", err)
	}
	if len(snapshot.Entries) != 0 {
		t.Fatalf("expected empty entries on permission error, got %#v", snapshot.Entries)
	}
}

func TestNormalizeLocalAddress(t *testing.T) {
	tests := map[string]struct {
		addr    string
		want    string
		include bool
		warn    bool
	}{
		"empty":       {addr: "", want: "0.0.0.0", include: true, warn: false},
		"loopback":    {addr: "127.0.0.1", want: "127.0.0.1", include: true, warn: false},
		"unspecified": {addr: "0.0.0.0", want: "0.0.0.0", include: true, warn: false},
		"ipv6loop":    {addr: "::1", want: "::1", include: true, warn: false},
		"ipv6unspec":  {addr: "[::]", want: "0.0.0.0", include: true, warn: false},
		"ipv6other":   {addr: "fe80::1", include: false, warn: true},
		"other":       {addr: "10.0.0.1", include: false, warn: false},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			gotAddr, gotInclude, gotWarn := normalizeLocalAddress(tt.addr)
			if gotAddr != tt.want {
				t.Fatalf("normalizeLocalAddress(%q) addr = %q, want %q", tt.addr, gotAddr, tt.want)
			}
			if gotInclude != tt.include {
				t.Fatalf(
					"normalizeLocalAddress(%q) include = %v, want %v",
					tt.addr,
					gotInclude,
					tt.include,
				)
			}
			if gotWarn != tt.warn {
				t.Fatalf("normalizeLocalAddress(%q) warn = %v, want %v", tt.addr, gotWarn, tt.warn)
			}
		})
	}
}
