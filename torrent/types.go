package torrent

type MetaInfo struct {
	AnnounceURL string
	Info        TorrentInfo
}

type TorrentInfo struct {
	PieceLength uint
	Pieces      string
	Files       []File
}

type File struct {
	Length uint
	Path   string
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
