package main

import (
	"fmt"
	"morpho/bencoding"
	"morpho/pwp"
	"morpho/torrent"
	"os"
	"sync"
	"time"
)

var (
	announceData torrent.AnnounceData
	Meta_info    torrent.MetaInfo
)

func init() {
	data, _ := os.ReadFile("another.torrent")
	source := string(data)
	bval, _ := bencoding.Decode(source)
	m := bval.(map[string]any)
	Meta_info, _ = torrent.LoadTorrent(bval)
	for index, hash := range Meta_info.Info.Pieces {
		fmt.Printf("HASH: %x INDEX: %d\n", hash, index)

	}
	fmt.Println(Meta_info.Info.PieceLength, Meta_info.Info.Files[0].Length)
	announceData = torrent.CreateAnnounceData(&Meta_info, m)

}

func main() {
	var peerList []torrent.Peer
	var wg sync.WaitGroup

	announceData.ManageAnnounceTracker(&Meta_info, &peerList)
	wg.Add(1)
	go func() {
		time.Sleep(20 * time.Second)
		pwp.PeerManagers.HandlePeer(&Meta_info)
		defer wg.Done()

	}()
	wg.Add(1)
	go pwp.StartPeerManager(&peerList, &announceData)
	wg.Wait()

	fmt.Println("finishing")
}
