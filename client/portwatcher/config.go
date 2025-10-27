package portwatcher

import "strings"

// Config describes filtering options applied to detected listening sockets.
type Config struct {
	MinPort          uint16
	MaxPort          uint16
	ExcludePorts     []uint16
	ExcludeProcesses []string
}

// ShouldInclude judges whether the provided entry satisfies the configuration.
func (c Config) ShouldInclude(entry Entry) bool {
	if c.MinPort > 0 && entry.Port < c.MinPort {
		return false
	}
	if c.MaxPort > 0 && entry.Port > c.MaxPort {
		return false
	}
	for _, port := range c.ExcludePorts {
		if entry.Port == port {
			return false
		}
	}
	if len(c.ExcludeProcesses) > 0 && entry.Process != "" {
		entryProcess := strings.ToLower(entry.Process)
		for _, process := range c.ExcludeProcesses {
			if entryProcess == strings.ToLower(process) {
				return false
			}
		}
	}
	return true
}
