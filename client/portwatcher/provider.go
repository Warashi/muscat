package portwatcher

import (
	"context"
	"errors"
	"log"
	"net/netip"
	"os"
	"strings"

	gnet "github.com/shirou/gopsutil/v4/net"
	gprocess "github.com/shirou/gopsutil/v4/process"
)

type systemSnapshotProvider struct {
	cfg             Config
	fetcher         func(context.Context) ([]gnet.ConnectionStat, error)
	processResolver func(context.Context, int32) (string, error)
	logger          *log.Logger
}

func newSystemSnapshotProvider(cfg Config, logger *log.Logger) snapshotProvider {
	return &systemSnapshotProvider{
		cfg: cfg,
		fetcher: func(ctx context.Context) ([]gnet.ConnectionStat, error) {
			return gnet.ConnectionsWithContext(ctx, "tcp")
		},
		processResolver: resolveProcessName,
		logger:          logger,
	}
}

func (p *systemSnapshotProvider) Snapshot(ctx context.Context) (Snapshot, error) {
	conns, err := p.fetcher(ctx)
	if err != nil {
		if errors.Is(err, os.ErrPermission) {
			p.logger.Printf("portwatcher: permission denied while enumerating connections: %v", err)
			return Snapshot{}, nil
		}
		return Snapshot{}, err
	}

	entries := make([]Entry, 0, len(conns))
	processCache := make(map[int32]string)
	for _, conn := range conns {
		if conn.Status != "LISTEN" {
			continue
		}

		addr, include, warnIPv6 := normalizeLocalAddress(conn.Laddr.IP)
		if warnIPv6 {
			p.logger.Printf(
				"portwatcher: skipping unsupported IPv6 address %s:%d",
				conn.Laddr.IP,
				conn.Laddr.Port,
			)
		}
		if !include {
			continue
		}

		entry := Entry{
			Addr:     addr,
			Port:     uint16(conn.Laddr.Port),
			Protocol: "tcp",
		}

		if !p.cfg.ShouldInclude(entry) {
			continue
		}

		if cached, ok := processCache[conn.Pid]; ok {
			entry.Process = cached
		} else if conn.Pid > 0 {
			name, err := p.processResolver(ctx, conn.Pid)
			if err != nil {
				if errors.Is(err, os.ErrPermission) {
					p.logger.Printf("portwatcher: permission denied while resolving process name for pid %d", conn.Pid)
				} else {
					p.logger.Printf("portwatcher: failed to resolve process name for pid %d: %v", conn.Pid, err)
				}
			} else if name != "" {
				processCache[conn.Pid] = name
				entry.Process = name
			}
		}

		if !p.cfg.ShouldInclude(entry) {
			continue
		}

		entries = append(entries, entry)
	}

	return Snapshot{Entries: entries}, nil
}

func resolveProcessName(ctx context.Context, pid int32) (string, error) {
	proc, err := gprocess.NewProcess(pid)
	if err != nil {
		return "", err
	}
	name, err := proc.NameWithContext(ctx)
	if err != nil {
		return "", err
	}
	return name, nil
}

func normalizeLocalAddress(raw string) (addr string, include bool, warnIPv6 bool) {
	trimmed := strings.Trim(raw, "[]")
	if trimmed == "" {
		return "0.0.0.0", true, false
	}

	parsed, err := netip.ParseAddr(trimmed)
	if err != nil {
		return "", false, false
	}

	unmapped := parsed.Unmap()
	if unmapped.Is4() {
		if unmapped.IsLoopback() || unmapped == netip.IPv4Unspecified() {
			return unmapped.String(), true, false
		}
		return "", false, false
	}

	if parsed == netip.IPv6Loopback() {
		return "::1", true, false
	}
	if parsed == netip.IPv6Unspecified() {
		return "0.0.0.0", true, false
	}
	return "", false, true
}
