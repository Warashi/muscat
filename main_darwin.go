//go:build darwin

package main

import "github.com/Warashi/go-swim"

func main() {
	ch := _main()
	for {
		select {
		case <-ch:
			return
		default:
			swim.EventLoop()
		}
	}
}
