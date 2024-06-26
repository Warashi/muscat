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
	"io"
	"log"
	"os"

	"github.com/spf13/cobra"

	"github.com/Warashi/muscat/v2/client"
)

// pasteCmd represents the paste command
var pasteCmd = &cobra.Command{
	Use:   "paste",
	Short: "Output clipboard contents to stdout",
	Long:  `Output clipboard contents of server host to stdout.`,
	Run: func(cmd *cobra.Command, args []string) {
		muscat := client.New(mustGetListenArgs(context.Background()))
		r, err := muscat.Paste(context.Background())
		if err != nil {
			log.Fatalf("muscat.Paste: %v", err)
		}
		if _, err := io.Copy(os.Stdout, r); err != nil {
			log.Fatalf("io.Copy: %v", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(pasteCmd)
}
