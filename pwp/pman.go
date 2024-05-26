package pwp

import (
	"fmt"
	"morpho/torrent"
	"morpho/utp"
	"strconv"
	"sync"
	"time"
)

var PeerManagers PeerManager
var Msg chan int

func (p *PeerManager) HandlePeer(metaInfo *torrent.MetaInfo) {
	// TODO else make request messages
	// TODO
	// each peer have go routines to handle incoming traffic
	for index, peer := range p.peers {
		fmt.Println("peers are :", peer.bitfield, index)

		peer.PwpMessage(Unchoke)
		peer.PwpMessage(Interested)
		peer.RequestMessage(0, metaInfo)

		res_buf := make([]byte, 10000)
		n, err := peer.utp.BaseConn.Read(res_buf)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(res_buf[:n])
		peer.HandleIncomingMessage(res_buf[:n])

	}

}

func (p *PeerInfo) HandleIncomingMessage(msg []byte) {
	packet := p.utp.ProcessReceivedPacket(msg)
	fmt.Println(packet)
	if len(packet.Payload) > 0 {

		msgId := packet.Payload[4]
		switch Id(msgId) {
		case Choke:
			p.choked = true
		case Interested:
			p.interested = true
		case Peice:
			fmt.Println(packet.Payload)
			// parse teh peice
		default:
			fmt.Println(packet.Payload)

		}
	}

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
				peer.utp = connection
				PeerManagers.peers = append(PeerManagers.peers, peer)
				fmt.Println(len(PeerManagers.peers))
				done <- true
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
