package main

import (
	"fmt"
	"morpho/bencoding"
	"morpho/torrent"
	"os"
)

func main() {
	data, _ := os.ReadFile("another.torrent")
	source := string(data)
	bval, _ := bencoding.Decode(source)
	m := bval.(map[string]any)
	meta_info, _ := torrent.LoadTorrent(bval)
	announceData := torrent.CreateAnnounceData(&meta_info, m)
	body := announceData.ManageAnnounceTracker(&meta_info)
	fmt.Println("--------this is main---------", body)
}
