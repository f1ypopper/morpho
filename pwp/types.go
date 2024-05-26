package pwp

import (
	"morpho/utp"
)

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
	peers      []PeerInfo
	Downloaded []byte
	left       []byte
}
type PeerInfo struct {
	utp        utp.UTPConnection
	bitfield   []byte
	interested bool
	choked     bool
	busy       bool
	done       chan bool
}
