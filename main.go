package main

import (
	"fmt"
	"morpho/bencoding"
	"morpho/pwp"
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
	for index, hash := range meta_info.Info.Pieces {
		fmt.Printf("HASH: %x INDEX: %d\n", hash, index)

	}
	fmt.Println(meta_info.Info.PieceLength, meta_info.Info.Files[0].Length)
	announceData = torrent.CreateAnnounceData(&meta_info, m)

}

func main() {
	// fmt.Println("ANNOUNCE LIST: ", meta_info.AnnounceList)
	// fmt.Printf("INFO HASH: %x\n", announceData.InfoHash)
	var peerList []torrent.Peer

	announceData.ManageAnnounceTracker(&meta_info, &peerList)
	fmt.Println("reading trackers is complete")
	pwp.StartPeerManager(&peerList, &announceData)
	fmt.Println("finishing")
}
