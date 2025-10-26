package cmd

import (
	"bytes"
	"context"
	"testing"

	"github.com/spf13/cobra"

	clientexposeport "github.com/Warashi/muscat/v2/client/exposeport"
)

func TestParsePortSpec(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		input       string
		wantLocal   uint16
		wantRemote  uint16
		expectError bool
	}{
		"local only":      {"8080", 8080, 8080, false},
		"local remote":    {"8080:9090", 8080, 9090, false},
		"remote zero":     {"8080:0", 8080, 0, false},
		"invalid format":  {"8080:9090:1", 0, 0, true},
		"missing remote":  {"8080:", 0, 0, true},
		"non numeric":     {"abc", 0, 0, true},
		"remote invalid":  {"8080:xyz", 0, 0, true},
		"local zero":      {"0:1234", 0, 0, true},
		"remote negative": {"8080:-1", 0, 0, true},
	}

	for name, tt := range tests {
		tt := tt
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			local, remote, err := parsePortSpec(tt.input)
			if tt.expectError {
				if err == nil {
					t.Fatalf("expected error for %q", tt.input)
				}
				return
			}
			if err != nil {
				t.Fatalf("parsePortSpec(%q): %v", tt.input, err)
			}
			if local != tt.wantLocal || remote != tt.wantRemote {
				t.Fatalf(
					"parsePortSpec(%q): got (%d,%d) want (%d,%d)",
					tt.input,
					local,
					remote,
					tt.wantLocal,
					tt.wantRemote,
				)
			}
		})
	}
}

func TestExposePortManual(t *testing.T) {
	t.Parallel()

	fakeMgr := &fakeExposePortManager{}
	var capturedOpts exposePortOptions
	deps := exposePortDeps{
		newManager: func(ctx context.Context, cmd *cobra.Command, opts exposePortOptions) (exposePortManager, error) {
			capturedOpts = opts
			return fakeMgr, nil
		},
		wait: func(context.Context) error { return nil },
	}

	command := newExposePortCmd(deps)
	command.SetContext(context.Background())
	var stdout, stderr bytes.Buffer
	command.SetOut(&stdout)
	command.SetErr(&stderr)
	command.SetArgs([]string{"--remote-policy", "next-free", "8080", "9090:9000"})

	if err := command.Execute(); err != nil {
		t.Fatalf("command.Execute: %v", err)
	}

	if len(fakeMgr.exposures) != 2 {
		t.Fatalf("expected 2 exposures, got %d", len(fakeMgr.exposures))
	}
	if fakeMgr.exposures[0] != (portSpec{local: 8080, remote: 8080}) {
		t.Fatalf("unexpected first exposure: %#v", fakeMgr.exposures[0])
	}
	if fakeMgr.exposures[1] != (portSpec{local: 9090, remote: 9000}) {
		t.Fatalf("unexpected second exposure: %#v", fakeMgr.exposures[1])
	}
	if capturedOpts.remotePolicy != clientexposeport.RemotePortPolicyNextFree {
		t.Fatalf("expected remote policy next-free, got %v", capturedOpts.remotePolicy)
	}
	if capturedOpts.bindAddress != "127.0.0.1" {
		t.Fatalf("expected default bind address, got %q", capturedOpts.bindAddress)
	}
}

func TestExposePortRequiresPorts(t *testing.T) {
	t.Parallel()

	deps := exposePortDeps{
		newManager: func(context.Context, *cobra.Command, exposePortOptions) (exposePortManager, error) {
			return &fakeExposePortManager{}, nil
		},
		wait: func(context.Context) error { return nil },
	}

	command := newExposePortCmd(deps)
	command.SetContext(context.Background())
	command.SetArgs([]string{})

	err := command.Execute()
	if err == nil || err.Error() != "at least one port must be specified" {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestExposePortAutoNotImplemented(t *testing.T) {
	t.Parallel()

	deps := exposePortDeps{
		newManager: func(context.Context, *cobra.Command, exposePortOptions) (exposePortManager, error) {
			return &fakeExposePortManager{}, nil
		},
		wait: func(context.Context) error { return nil },
	}

	command := newExposePortCmd(deps)
	command.SetContext(context.Background())
	command.SetArgs([]string{"--auto"})

	err := command.Execute()
	if err == nil || err.Error() != "--auto is not implemented yet" {
		t.Fatalf("unexpected error: %v", err)
	}
}

type fakeExposePortManager struct {
	exposures []portSpec
}

func (m *fakeExposePortManager) Expose(ctx context.Context, local, remote uint16) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}
	m.exposures = append(m.exposures, portSpec{local: local, remote: remote})
	return nil
}

func (m *fakeExposePortManager) Stop(uint16) {}

func (m *fakeExposePortManager) Shutdown() {}
