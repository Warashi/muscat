package portwatcher

import (
	"sort"
	"strings"
)

// Entry represents a single listening socket on the local machine.
type Entry struct {
	Port     uint16
	Addr     string
	Protocol string
	Process  string
}

// Snapshot holds a collection of listening socket entries at a given point in time.
type Snapshot struct {
	Entries []Entry
}

// Diff captures the difference between two snapshots.
type Diff struct {
	Added   []Entry
	Removed []Entry
}

// Empty reports whether the diff does not contain any changes.
func (d Diff) Empty() bool {
	return len(d.Added) == 0 && len(d.Removed) == 0
}

// Diff calculates differences compared with the provided previous snapshot.
func (s Snapshot) Diff(prev Snapshot) Diff {
	currentEntries := deduplicateAndSort(s.Entries)
	previousEntries := deduplicateAndSort(prev.Entries)

	current := make(map[string]Entry, len(currentEntries))
	for _, entry := range currentEntries {
		current[entryKey(entry)] = entry
	}

	previous := make(map[string]Entry, len(previousEntries))
	for _, entry := range previousEntries {
		previous[entryKey(entry)] = entry
	}

	var diff Diff
	for key, entry := range current {
		if _, ok := previous[key]; !ok {
			diff.Added = append(diff.Added, entry)
		}
	}
	for key, entry := range previous {
		if _, ok := current[key]; !ok {
			diff.Removed = append(diff.Removed, entry)
		}
	}

	sort.Slice(diff.Added, func(i, j int) bool {
		return entryKey(diff.Added[i]) < entryKey(diff.Added[j])
	})
	sort.Slice(diff.Removed, func(i, j int) bool {
		return entryKey(diff.Removed[i]) < entryKey(diff.Removed[j])
	})

	return diff
}

func deduplicateAndSort(entries []Entry) []Entry {
	if len(entries) == 0 {
		return nil
	}
	copied := append([]Entry(nil), entries...)
	sort.Slice(copied, func(i, j int) bool {
		return entryKey(copied[i]) < entryKey(copied[j])
	})
	deduplicated := copied[:1]
	for i := 1; i < len(copied); i++ {
		if entryKey(copied[i]) == entryKey(copied[i-1]) {
			continue
		}
		deduplicated = append(deduplicated, copied[i])
	}
	return deduplicated
}

func entryKey(entry Entry) string {
	var builder strings.Builder
	builder.Grow(len(entry.Addr) + len(entry.Protocol) + len(entry.Process) + 10)
	builder.WriteString(entry.Addr)
	builder.WriteString("|")
	builder.WriteString(entry.Protocol)
	builder.WriteString("|")
	builder.WriteString(entry.Process)
	builder.WriteString("|")
	builder.WriteString(portKey(entry.Port))
	return builder.String()
}

func portKey(port uint16) string {
	// Value is at most 5 digits, but reserve some padding for safety.
	var buf [6]byte
	i := len(buf)
	for {
		i--
		buf[i] = byte('0' + port%10)
		port /= 10
		if port == 0 {
			break
		}
	}
	return string(buf[i:])
}
