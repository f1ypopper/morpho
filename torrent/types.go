package torrent

import (
	"net"
	"net/url"
)

type MetaInfo struct {
	AnnounceURL  string
	AnnounceList []url.URL
	Info         TorrentInfo
	Length       uint
}

type TorrentInfo struct {
	PieceLength uint
	Pieces      string
	Files       []File
}

type File struct {
	Length uint
	Path   []any
}

type Peer struct {
	IP     net.IP
	PeerID string
	Port   uint16
}

type ResponseData struct {
	Complete    uint
	Incomplete  uint
	Interval    uint
	MinInterval uint
	Peer        []Peer
}

type AnnounceData struct {
	InfoHash      []byte
	PeerID        string
	Port          uint16
	Uploaded      uint64
	Downloaded    uint64
	Left          uint64
	Compact       bool
	Event         string
	ConnectionID  uint64
	Action        uint32
	TransactionID uint32
}
