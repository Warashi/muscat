package cmd

import (
	"bytes"
	"context"
	"testing"
	"time"

	"github.com/spf13/cobra"

	clientexposeport "github.com/Warashi/muscat/v2/client/exposeport"
	"github.com/Warashi/muscat/v2/client/portwatcher"
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
	if capturedOpts.autoInterval != portwatcher.DefaultInterval {
		t.Fatalf("expected default auto interval, got %v", capturedOpts.autoInterval)
	}
	if len(capturedOpts.autoExclude) != 0 {
		t.Fatalf("expected no auto exclude ports, got %v", capturedOpts.autoExclude)
	}
	if len(capturedOpts.autoExcludeProcess) != 0 {
		t.Fatalf("expected no auto exclude processes, got %v", capturedOpts.autoExcludeProcess)
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

func TestExposePortAuto(t *testing.T) {
	t.Parallel()

	fakeMgr := &fakeExposePortManager{}
	var (
		started   bool
		gotOpts   exposePortOptions
		gotManual map[uint16]struct{}
	)

	deps := exposePortDeps{
		newManager: func(context.Context, *cobra.Command, exposePortOptions) (exposePortManager, error) {
			return fakeMgr, nil
		},
		wait: func(context.Context) error { return nil },
		startAuto: func(ctx context.Context, cmd *cobra.Command, opts exposePortOptions, manager exposePortManager, manual map[uint16]struct{}) error {
			started = true
			gotOpts = opts
			copied := make(map[uint16]struct{}, len(manual))
			for k, v := range manual {
				copied[k] = v
			}
			gotManual = copied
			return nil
		},
	}

	command := newExposePortCmd(deps)
	command.SetContext(context.Background())
	var stdout, stderr bytes.Buffer
	command.SetOut(&stdout)
	command.SetErr(&stderr)
	command.SetArgs([]string{
		"--auto",
		"--auto-interval", "2s",
		"--auto-exclude", "1000",
		"--auto-exclude", "2000",
		"--auto-exclude-process", "foo",
		"--auto-exclude-process", "bar",
		"5000",
	})

	if err := command.Execute(); err != nil {
		t.Fatalf("command.Execute: %v", err)
	}

	if !started {
		t.Fatalf("startAuto was not called")
	}
	if gotOpts.autoInterval != 2*time.Second {
		t.Fatalf("unexpected auto interval: %v", gotOpts.autoInterval)
	}
	if len(gotOpts.autoExclude) != 2 || gotOpts.autoExclude[0] != 1000 ||
		gotOpts.autoExclude[1] != 2000 {
		t.Fatalf("unexpected auto exclude ports: %v", gotOpts.autoExclude)
	}
	if len(gotOpts.autoExcludeProcess) != 2 || gotOpts.autoExcludeProcess[0] != "foo" ||
		gotOpts.autoExcludeProcess[1] != "bar" {
		t.Fatalf("unexpected auto exclude processes: %v", gotOpts.autoExcludeProcess)
	}
	if _, ok := gotManual[5000]; !ok {
		t.Fatalf("manual port not propagated to auto starter")
	}
	if len(fakeMgr.exposures) != 1 || fakeMgr.exposures[0].local != 5000 {
		t.Fatalf("manual exposure not started: %#v", fakeMgr.exposures)
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
