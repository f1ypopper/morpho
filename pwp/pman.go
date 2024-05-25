package pwp

import (
	"fmt"
	"morpho/torrent"
	"morpho/utp"
	"strconv"
	"sync"
	"time"
)

func HandlePeer() {
	// TODO Start handshake
	// TODO create channel with peer manager and data manager
	// TODO if bitfield contact peer manager
	// TODO else make request messages
	// TODO
}

func (pinfo *PeerInfo) timer() {
	time.Sleep(time.Second * 10)
	pinfo.done <- true
	fmt.Println("timer", pinfo.done)

}

// TODO
// want handshake here utp.conn.MakePacket
// make a list of active peers map[ip]net.Conn

func StartPeerManager(pList *[]torrent.Peer, aData *torrent.AnnounceData) {

	var wg sync.WaitGroup
	handshakeData := handshake(aData)
	// ctx, cancel := context.WithCancel(context.Background())
	// go func() {
	// 	time.Sleep(time.Second * 10)
	// 	cancel()

	// }()
	for _, ipStr := range *pList {
		ipAdd := ipStr.IP.String() + ":" + strconv.Itoa(int(ipStr.Port))
		fmt.Println("conntecting to ", ipAdd)
		wg.Add(1)
		go func() {
			defer wg.Done()
			utp.Send(ipAdd, handshakeData)
		}()
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
