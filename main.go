package main

import (
	"os"
)

type Beencode struct {
}

func main() {
	if len(os.Args) < 2 {
		panic("Usage: torrent <torrent file>")
	}
	torrentPath := os.Args[1]
	if torrentPath == "" {
		panic("Usage: torrent <torrent file>")
	}
	println(torrentPath)
}
