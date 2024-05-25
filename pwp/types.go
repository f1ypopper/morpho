package pwp

import "net"

type HandshakeSt struct {
	Pstrlen  uint8
	Pstr     string
	Reserved []byte
	InfoHash []byte
	PeerId   []byte
}
type PieceInfo struct {
	state uint8
	// 0 available to download
	// 1 downloaded
	// 2 inprogress

}
type PeerManager struct {
}
type PeerInfo struct {
	ip              string
	conn            net.Conn
	interested      bool
	availablePieces uint
	busy            bool
	done            chan bool
}
