/*
Copyright © 2022 Shinnosuke Sawada <6warashi9@gmail.com>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package cmd

import (
	"context"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/spf13/cobra"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"

	"github.com/Warashi/muscat/v2/client"
	"github.com/Warashi/muscat/v2/pb/pbconnect"
	"github.com/Warashi/muscat/v2/server"
)

// serverCmd represents the server command
var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Start server for communicate with remote machine",
	Long:  `Start rpc server for communicate with client invoked at remote machine`,
	Run: func(cmd *cobra.Command, args []string) {
		network, addr := mustGetListenArgs(context.Background())
		muscatClient := client.New(network, addr)
		if network == "unix" {
			defer os.Remove(addr)
		}

		l, err := net.Listen(network, addr)
		if err != nil {
			if network == "unix" {
				// Remove addr file if already exists and retry
				if err := os.Remove(addr); err != nil {
					cmd.PrintErrf("os.Remove: %v", err)
					return
				}
				l, err = net.Listen(network, addr)
				if err != nil {
					cmd.PrintErrf("net.Listen: %v", err)
					return
				}
			}
		}

		mux := http.NewServeMux()
		mux.Handle(pbconnect.NewMuscatServiceHandler(new(server.MuscatServer)))

		go checkSocketProcess(muscatClient)

		if err := http.Serve(l, h2c.NewHandler(mux, new(http2.Server))); err != nil {
			cmd.PrintErrf("s.Serve: %v", err)
			return
		}
	},
}

func init() {
	rootCmd.AddCommand(serverCmd)
}

func getSocketProcessID(muscat *client.MuscatClient) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	pid, err := muscat.Health(ctx)
	if err != nil {
		return 0, err
	}
	return pid, nil
}

func checkSocketProcess(muscatClient *client.MuscatClient) {
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		pid, err := getSocketProcessID(muscatClient)
		if err != nil {
			continue
		}
		if pid != os.Getpid() {
			os.Exit(0)
		}
	}
}
