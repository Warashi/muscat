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
package main

import (
	"os"

	"github.com/Warashi/muscat/cmd"
)

func _main() chan struct{} {
	ch := make(chan struct{})
	go func() {
		defer close(ch)
		switch os.Args[0] {
		case "xdg-open":
			os.Args = append([]string{"muscat", "open"}, os.Args[1:]...)
		case "pbcopy":
			os.Args = append([]string{"muscat", "copy"}, os.Args[1:]...)
		case "pbpaste":
			os.Args = append([]string{"muscat", "paste"}, os.Args[1:]...)
		}
		cmd.Execute()
	}()
	return ch
}
