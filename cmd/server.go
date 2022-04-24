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
	"log"
	"net"
	"os"

	"github.com/spf13/cobra"
	"google.golang.org/grpc"

	"github.com/Warashi/muscat/pb"
	"github.com/Warashi/muscat/server"
)

// serverCmd represents the server command
var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Start server for communicate with remote machine",
	Long:  `Start rpc server for communicate with client invoked at remote machine`,
	Run: func(cmd *cobra.Command, args []string) {
		_ = os.Remove(mustGetSocketPath())
		l, err := net.Listen("unix", mustGetSocketPath())
		if err != nil {
			log.Fatalf("net.Listen: %v", err)
		}

		s := grpc.NewServer()
		pb.RegisterMuscatServer(s, new(server.Muscat))
		if err := s.Serve(l); err != nil {
			log.Fatalf("s.Serve: %v", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(serverCmd)
}