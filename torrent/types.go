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
	Port          uint
	Uploaded      uint
	Downloaded    uint
	Left          uint
	Compact       bool
	Event         string
	ConnectionID  uint64
	Action        uint32
	TransactionID uint32
}
