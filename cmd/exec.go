package cmd

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/spf13/cobra"

	"github.com/Warashi/muscat/v2/client"
)

// execCmd represents the exec command
var execCmd = &cobra.Command{
	Use:   "exec <command> [args...]",
	Short: "Execute a command on the server",
	Long: `Execute a command on the remote server and return the output.
The command can receive stdin input from the client.

Examples:
  # Simple command execution
  muscat exec ls -la

  # Command with stdin
  echo "Hello, World!" | muscat exec cat

  # Command with arguments
  muscat exec grep "pattern" file.txt`,
	Args: cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := context.Background()
		network, addr := mustGetListenArgs(ctx)

		// Read stdin if available
		var stdin []byte
		stat, _ := os.Stdin.Stat()
		if (stat.Mode() & os.ModeCharDevice) == 0 {
			// stdin is a pipe
			var err error
			stdin, err = io.ReadAll(os.Stdin)
			if err != nil {
				return fmt.Errorf("reading stdin: %w", err)
			}
		}

		// Create client and execute command
		c := client.New(network, addr)
		stdout, stderr, exitCode, err := c.Exec(ctx, args[0], args[1:], stdin)
		if err != nil {
			return fmt.Errorf("executing command: %w", err)
		}

		// Output results
		if len(stdout) > 0 {
			os.Stdout.Write(stdout)
		}
		if len(stderr) > 0 {
			os.Stderr.Write(stderr)
		}

		// Exit with the same code as the remote command
		if exitCode != 0 {
			os.Exit(exitCode)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(execCmd)
}
