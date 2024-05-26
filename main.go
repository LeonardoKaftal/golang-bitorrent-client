package main

import (
	"log"
	"main/torrent"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("MISSING TORRENT FILE PATH, Usage: torrent <torrent file path>")
	}
	torrentPath := os.Args[1]
	torrentFile, err := torrent.ReadTorrentFile(torrentPath)
	if err != nil {
		log.Fatal(err)
	}
	print(torrentFile)
}
