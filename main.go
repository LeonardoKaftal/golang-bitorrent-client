package main

import (
	"fmt"
	"log"
	"main/bencode_parser"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("Usage: torrent <torrent file>")
	}
	torrentPath := os.Args[1]
	torrentBuff, err := os.ReadFile(torrentPath)
	if err != nil {
		log.Fatal("Error reading torrent file:", err)
	}
	res := bencode_parser.Marshall(torrentBuff)
	fmt.Println(res)
}
