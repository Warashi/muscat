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
	"net"
	"net/http"
	"os"

	"github.com/spf13/cobra"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"

	"github.com/Warashi/muscat/v2/pb/pbconnect"
	"github.com/Warashi/muscat/v2/server"
)

// serverCmd represents the server command
var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Start server for communicate with remote machine",
	Long:  `Start rpc server for communicate with client invoked at remote machine`,
	Run: func(cmd *cobra.Command, args []string) {
		defer os.Remove(mustGetSocketPath())

		l, err := net.Listen("unix", mustGetSocketPath())
		if err != nil {
			cmd.PrintErrf("remove %s and try again\n", mustGetSocketPath())
			cmd.PrintErrf("net.Listen: %v", err)
			return
		}

		mux := http.NewServeMux()
		mux.Handle(pbconnect.NewMuscatServiceHandler(new(server.MuscatServer)))

		if err := http.Serve(l, h2c.NewHandler(mux, new(http2.Server))); err != nil {
			cmd.PrintErrf("s.Serve: %v", err)
			return
		}
	},
}

func init() {
	rootCmd.AddCommand(serverCmd)
}
