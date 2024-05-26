package pwp

import "morpho/torrent"

// func request message
// func choke and unchoke as well handle it
// request: <len=0013><id=6><index><begin><length>

// The request message is fixed length, and is used to request a block. The payload contains the following information:

//     index: integer specifying the zero-based piece index
//     begin: integer specifying the zero-based byte offset within the piece
//     length: integer specifying the requested length.

// availablePiece
const (
	Choke = iota
	Unchoke
	Interested
	Notinterested
	Have
	Bitfield
	Request
	Peice
	Cancel
)

type Message struct {
	ID      int
	Length  int
	Payload []byte
}

func (m *Message) RequestMessage(index int, metaInfo *torrent.MetaInfo) {
	peiceSize := metaInfo.Info.PieceLength
	totalsize := metaInfo.Info.Files[0].Length
	MaxPeiceSize := 16384 // 16 KB
	begin := index * int(peiceSize)
	end := begin + int(totalsize)
	for i := 0; i < end/MaxPeiceSize; i++ {
		// make request

		begin += MaxPeiceSize
		end += MaxPeiceSize

	}
}
