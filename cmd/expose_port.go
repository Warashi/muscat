package cmd

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/spf13/cobra"

	"github.com/Warashi/muscat/v2/client"
	clientexposeport "github.com/Warashi/muscat/v2/client/exposeport"
	"github.com/Warashi/muscat/v2/client/portwatcher"
)

type exposePortManager interface {
	Expose(context.Context, uint16, uint16) error
	Stop(uint16)
	Shutdown()
}

type exposePortDeps struct {
	newManager func(context.Context, *cobra.Command, exposePortOptions) (exposePortManager, error)
	wait       func(context.Context) error
	startAuto  func(context.Context, *cobra.Command, exposePortOptions, exposePortManager, map[uint16]struct{}) error
}

type exposePortOptions struct {
	bindAddress        string
	public             bool
	remotePolicy       clientexposeport.RemotePortPolicy
	auto               bool
	manual             []portSpec
	autoInterval       time.Duration
	autoExclude        []uint16
	autoExcludeProcess []string
}

type portSpec struct {
	local  uint16
	remote uint16
}

var defaultExposePortDeps = exposePortDeps{
	newManager: newDefaultExposePortManager,
	wait: func(ctx context.Context) error {
		<-ctx.Done()
		return ctx.Err()
	},
	startAuto: startAutoExpose,
}

// exposePortCmd represents the expose-port command.
func newExposePortCmd(deps exposePortDeps) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "expose-port [local[:remote]]...",
		Short: "Expose local listening ports to the remote server",
		Long: `Expose local listening ports to the remote server.

Specify ports as local[:remote]. When the remote port is omitted, the same number
as the local port is requested.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			opts, err := parseExposePortOptions(cmd, args)
			if err != nil {
				return err
			}
			if len(opts.manual) == 0 && !opts.auto {
				return errors.New("at least one port must be specified")
			}
			ctx, cancel := signal.NotifyContext(cmd.Context(), os.Interrupt, syscall.SIGTERM)
			defer cancel()

			manager, err := deps.newManager(ctx, cmd, opts)
			if err != nil {
				return err
			}
			defer manager.Shutdown()

			manualSet := make(map[uint16]struct{}, len(opts.manual))
			for _, spec := range opts.manual {
				if err := manager.Expose(ctx, spec.local, spec.remote); err != nil {
					return err
				}
				manualSet[spec.local] = struct{}{}
			}

			errCh := make(chan error, 2)
			go func() {
				errCh <- deps.wait(ctx)
			}()
			expected := 1

			if opts.auto {
				if deps.startAuto == nil {
					return errors.New("auto mode is not supported")
				}
				expected++
				go func() {
					errCh <- deps.startAuto(ctx, cmd, opts, manager, manualSet)
				}()
			}

			var firstErr error
			for i := 0; i < expected; i++ {
				if err := <-errCh; err != nil && !errors.Is(err, context.Canceled) &&
					firstErr == nil {
					firstErr = err
					cancel()
				}
			}
			if firstErr != nil {
				return firstErr
			}
			return nil
		},
		SilenceUsage:  true,
		SilenceErrors: false,
	}

	cmd.Flags().String("bind-address", "127.0.0.1", "Remote bind address for exposed ports")
	cmd.Flags().Bool("public", false, "Shortcut for bind-address=0.0.0.0")
	cmd.Flags().String("remote-policy", "fail", "Remote port conflict policy: fail|next-free|skip")
	cmd.Flags().Bool("auto", false, "Automatically expose newly opened local ports")
	cmd.Flags().
		Duration("auto-interval", portwatcher.DefaultInterval, "Polling interval used by --auto")
	cmd.Flags().UintSlice("auto-exclude", nil, "Ports to exclude from automatic exposure")
	cmd.Flags().
		StringSlice("auto-exclude-process", nil, "Process names to exclude from automatic exposure")

	return cmd
}

func parseExposePortOptions(cmd *cobra.Command, args []string) (exposePortOptions, error) {
	bindAddress, err := cmd.Flags().GetString("bind-address")
	if err != nil {
		return exposePortOptions{}, err
	}
	public, err := cmd.Flags().GetBool("public")
	if err != nil {
		return exposePortOptions{}, err
	}
	if public {
		bindAddress = "0.0.0.0"
	}

	policyText, err := cmd.Flags().GetString("remote-policy")
	if err != nil {
		return exposePortOptions{}, err
	}
	policy, err := clientexposeport.ParseRemotePortPolicy(policyText)
	if err != nil {
		return exposePortOptions{}, err
	}

	auto, err := cmd.Flags().GetBool("auto")
	if err != nil {
		return exposePortOptions{}, err
	}

	autoInterval, err := cmd.Flags().GetDuration("auto-interval")
	if err != nil {
		return exposePortOptions{}, err
	}
	if autoInterval <= 0 {
		autoInterval = portwatcher.DefaultInterval
	}

	autoExcludeRaw, err := cmd.Flags().GetUintSlice("auto-exclude")
	if err != nil {
		return exposePortOptions{}, err
	}
	autoExclude := make([]uint16, 0, len(autoExcludeRaw))
	for _, v := range autoExcludeRaw {
		if v > 65535 {
			return exposePortOptions{}, fmt.Errorf("auto exclude port %d out of range", v)
		}
		autoExclude = append(autoExclude, uint16(v))
	}

	autoExcludeProcess, err := cmd.Flags().GetStringSlice("auto-exclude-process")
	if err != nil {
		return exposePortOptions{}, err
	}

	manual, err := parsePortSpecs(args)
	if err != nil {
		return exposePortOptions{}, err
	}

	return exposePortOptions{
		bindAddress:        bindAddress,
		public:             public,
		remotePolicy:       policy,
		auto:               auto,
		manual:             manual,
		autoInterval:       autoInterval,
		autoExclude:        autoExclude,
		autoExcludeProcess: autoExcludeProcess,
	}, nil
}

func parsePortSpecs(args []string) ([]portSpec, error) {
	specs := make([]portSpec, 0, len(args))
	for _, arg := range args {
		arg = strings.TrimSpace(arg)
		if arg == "" {
			return nil, errors.New("port spec cannot be empty")
		}
		local, remote, err := parsePortSpec(arg)
		if err != nil {
			return nil, err
		}
		specs = append(specs, portSpec{local: local, remote: remote})
	}
	return specs, nil
}

func parsePortSpec(input string) (uint16, uint16, error) {
	parts := strings.Split(input, ":")
	if len(parts) > 2 {
		return 0, 0, fmt.Errorf("invalid port spec %q", input)
	}
	local, err := parsePortNumber(parts[0], false)
	if err != nil {
		return 0, 0, fmt.Errorf("invalid local port in %q: %w", input, err)
	}
	remote := local
	if len(parts) == 2 {
		if parts[1] == "" {
			return 0, 0, fmt.Errorf("remote port missing in %q", input)
		}
		remote, err = parsePortNumber(parts[1], true)
		if err != nil {
			return 0, 0, fmt.Errorf("invalid remote port in %q: %w", input, err)
		}
	}
	return local, remote, nil
}

func parsePortNumber(value string, allowZero bool) (uint16, error) {
	port, err := strconv.ParseUint(value, 10, 16)
	if err != nil {
		return 0, err
	}
	if port == 0 && !allowZero {
		return 0, errors.New("port must be greater than zero")
	}
	if port > 65535 {
		return 0, errors.New("port must be <= 65535")
	}
	return uint16(port), nil
}

func startAutoExpose(
	ctx context.Context,
	cmd *cobra.Command,
	opts exposePortOptions,
	manager exposePortManager,
	manual map[uint16]struct{},
) error {
	logger := log.New(cmd.ErrOrStderr(), "", log.LstdFlags)
	pwCfg := portwatcher.Config{
		ExcludePorts:     append([]uint16(nil), opts.autoExclude...),
		ExcludeProcesses: append([]string(nil), opts.autoExcludeProcess...),
	}
	watcherOpts := []portwatcher.Option{
		portwatcher.WithLogger(logger),
	}
	if opts.autoInterval > 0 {
		watcherOpts = append(watcherOpts, portwatcher.WithInterval(opts.autoInterval))
	}

	watcher := portwatcher.New(pwCfg, watcherOpts...)

	active := make(map[uint16]struct{})
	var mu sync.Mutex

	handler := func(diff portwatcher.Diff) {
		var toStart, toStop []uint16

		mu.Lock()
		for _, entry := range diff.Added {
			if _, ok := manual[entry.Port]; ok {
				continue
			}
			if _, ok := active[entry.Port]; ok {
				continue
			}
			active[entry.Port] = struct{}{}
			toStart = append(toStart, entry.Port)
		}
		for _, entry := range diff.Removed {
			if _, ok := manual[entry.Port]; ok {
				continue
			}
			if _, ok := active[entry.Port]; ok {
				delete(active, entry.Port)
				toStop = append(toStop, entry.Port)
			}
		}
		mu.Unlock()

		for _, port := range toStart {
			if err := manager.Expose(ctx, port, port); err != nil {
				mu.Lock()
				delete(active, port)
				mu.Unlock()
				fmt.Fprintf(cmd.ErrOrStderr(), "failed to auto expose port %d: %v\n", port, err)
				continue
			}
			fmt.Fprintf(cmd.OutOrStdout(), "auto exposing local %d\n", port)
		}
		for _, port := range toStop {
			manager.Stop(port)
			fmt.Fprintf(cmd.OutOrStdout(), "auto stopped local %d\n", port)
		}
	}

	if err := watcher.Start(ctx, handler); err != nil && !errors.Is(err, context.Canceled) {
		return err
	}
	return nil
}

func newDefaultExposePortManager(
	ctx context.Context,
	cmd *cobra.Command,
	opts exposePortOptions,
) (exposePortManager, error) {
	logger := log.New(cmd.ErrOrStderr(), "", log.LstdFlags)
	network, addr := mustGetListenArgs(ctx)
	muscat := client.New(network, addr)

	cfg := clientexposeport.ManagerConfig{
		BindAddress: opts.bindAddress,
		Public:      opts.public,
		Policy:      opts.remotePolicy,
		Logger:      logger,
		Handlers: clientexposeport.EventHandlers{
			Ready: func(binding clientexposeport.Binding) {
				fmt.Fprintf(
					cmd.OutOrStdout(),
					"local %d exposed as remote %d\n",
					binding.LocalPort,
					binding.RemotePort,
				)
			},
			Error: func(binding clientexposeport.Binding, err error) {
				fmt.Fprintf(
					cmd.ErrOrStderr(),
					"failed to expose local %d (remote %d): %v\n",
					binding.LocalPort,
					binding.RemotePort,
					err,
				)
			},
			Stopped: func(binding clientexposeport.Binding, err error) {
				if err != nil {
					fmt.Fprintf(
						cmd.ErrOrStderr(),
						"stopped exposing local %d (remote %d): %v\n",
						binding.LocalPort,
						binding.RemotePort,
						err,
					)
				}
			},
		},
	}
	return clientexposeport.NewManager(muscat, cfg), nil
}

func init() {
	rootCmd.AddCommand(newExposePortCmd(defaultExposePortDeps))
}
