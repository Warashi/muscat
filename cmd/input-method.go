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
	"fmt"
	"log"

	"github.com/spf13/cobra"

	"github.com/Warashi/muscat/client"
)

// getInputMethodCmd represents the get-im command
var getInputMethodCmd = &cobra.Command{
	Use:   "get-im",
	Short: "get input method of server host",
	Long:  `get input method of server host`,
	Run: func(cmd *cobra.Command, args []string) {
		muscat, err := client.New(mustGetSocketPath())
		if err != nil {
			log.Fatalf("client.New: %v", err)
		}
		id, err := muscat.GetInputMethod(context.Background())
		if err != nil {
			log.Fatalf("muscat.GetInputMethod: %v", err)
		}
		fmt.Println(id)
	},
}

// setInputMethodCmd represents the set-im command
var setInputMethodCmd = &cobra.Command{
	Use:   "set-im",
	Short: "set input method of server host",
	Long:  `set input method of server host`,
	Run: func(cmd *cobra.Command, args []string) {
		muscat, err := client.New(mustGetSocketPath())
		if err != nil {
			log.Fatalf("client.New: %v", err)
		}
		if err := muscat.SetInputMethod(context.Background(), args[0]); err != nil {
			log.Fatalf("muscat.SetInputMethod: %v", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(getInputMethodCmd)
	rootCmd.AddCommand(setInputMethodCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// copyCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// copyCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}