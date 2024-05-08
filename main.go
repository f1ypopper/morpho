package main

import (
	"fmt"
	"morpho/bencoding"
	"morpho/torrent"
	"morpho/utp"
	"os"
	"time"
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

//func main() {
//	fmt.Printf("INFO HASH: %x\n", announceData.InfoHash)
//	var peerList = map[string]uint16{}
//	announceData.ManageAnnounceTracker(&meta_info, &peerList)
//	pwp.StartPeerManager(&peerList, &announceData)
//}

func main() {
	conn, err := utp.Dial("localhost:1111")
	if err != nil {
		panic(err.Error())
	}
	conn.Write([]byte("Hello World\n"))
	buf := make([]byte, 10)
	conn.Read(buf)
	fmt.Printf("READ %d BYTES: %s\n", len(buf), buf)
	conn.Write([]byte("Bye World\n"))
	buf = make([]byte, 5)
	conn.Read(buf)
	fmt.Printf("READ %d BYTES: %s\n", len(buf), buf)
	time.Sleep(10000 * time.Second)
}
