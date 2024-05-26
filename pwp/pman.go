package pwp

import (
	"fmt"
	"morpho/torrent"
	"morpho/utp"
	"strconv"
	"sync"
	"time"
)

var peerManager PeerManager

func (p *PeerManager) HandlePeer() {
	// TODO Start handshake
	// TODO create channel with peer manager and data manager
	// TODO if bitfield contact peer manager
	// TODO else make request messages
	// TODO

}

// TODO
// want handshake here utp.conn.MakePacket
// make a list of active peers map[ip]net.Conn

func StartPeerManager(pList *[]torrent.Peer, aData *torrent.AnnounceData) {

	var wg sync.WaitGroup
	handshakeData := handshake(aData)
	done := make(chan bool)

	go func() {
		time.Sleep(time.Second * 10)
		done <- true
	}()
	fmt.Println("lenght of peer list is ", len(*pList))

	for _, ipStr := range *pList {
		var peer PeerInfo
		ipAdd := ipStr.IP.String() + ":" + strconv.Itoa(int(ipStr.Port))
		wg.Add(10)
		go func(ipAdd string) {
			defer wg.Done()
			bitfield, connection, err := utp.NewPeer(ipAdd, handshakeData)
			if err != nil {
				return

			}
			if bitfield != nil {
				peer.bitfield = bitfield
				peer.conn = connection
				peerManager.peers = append(peerManager.peers, peer)
				fmt.Println(len(peerManager.peers))
			}

		}(ipAdd)
	}

	wg.Wait()

}

func handshake(aD *torrent.AnnounceData) []byte {
	buf := make([]byte, 68)
	buf[0] = 19
	curr := 1
	curr += copy(buf[curr:], "BitTorrent protocol")
	curr += copy(buf[curr:], make([]byte, 8)) // 8 reserved bytes
	curr += copy(buf[curr:], aD.InfoHash[:])
	curr += copy(buf[curr:], aD.PeerID[:])
	return buf
}
