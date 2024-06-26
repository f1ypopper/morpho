package main

import (
	"fmt"
	"morpho/bencoding"
	"morpho/torrent"
	"os"
)

var (
	announceData torrent.AnnounceData
	meta_info    torrent.MetaInfo
)

func init() {
	data, _ := os.ReadFile("another.torrent")
	source := string(data)
	bval, _ := bencoding.Decode(source)
	m := bval.(map[string]any)
	meta_info, _ = torrent.LoadTorrent(bval)
	announceData = torrent.CreateAnnounceData(&meta_info, m)

}

func main() {
	fmt.Println("ANNOUNCE LIST: ", meta_info.AnnounceList)
	fmt.Printf("INFO HASH: %x\n", announceData.InfoHash)
	var peerList = map[string]uint16{}
	announceData.ManageAnnounceTracker(&meta_info, &peerList)
}
