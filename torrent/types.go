package torrent

type MetaInfo struct {
	Announce string
	Info     TorrentInfo
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
