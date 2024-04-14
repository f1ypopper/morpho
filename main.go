package main

import (
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
	var peerList = []torrent.Peer{}
	announceData.ManageAnnounceTracker(&meta_info, &peerList)
}
