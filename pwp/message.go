package pwp

import (
	"encoding/binary"
	"morpho/torrent"
)

// func request message
// func choke and unchoke as well handle it
// request: <len=0013><id=6><index><begin><length>

// The request message is fixed length, and is used to request a block. The payload contains the following information:

//     index: integer specifying the zero-based piece index
//     begin: integer specifying the zero-based byte offset within the piece
//     length: integer specifying the requested length.

// availablePiece
type Id int

const (
	Choke Id = iota
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
	ID      Id
	Length  int
	Payload []byte
}

func (p *PeerInfo) RequestMessage(index int, metaInfo *torrent.MetaInfo) {
	var m Message
	peiceSize := metaInfo.Info.PieceLength
	totalsize := metaInfo.Info.Files[0].Length
	MaxPeiceSize := 16384 // 16 KB
	begin := index * int(peiceSize)
	end := begin + int(totalsize)
	m.ID = Request
	m.Length = 13
	for i := 0; i < end/MaxPeiceSize; i++ {
		// craft payload message
		payload := make([]byte, 8)

		binary.BigEndian.PutUint32(payload[0:], uint32(index))
		binary.BigEndian.PutUint16(payload[4:], uint16(begin))
		binary.BigEndian.PutUint16(payload[4:], uint16(MaxPeiceSize))
		m.Payload = append(m.Payload, payload...)

		begin += MaxPeiceSize
		end += MaxPeiceSize

	}
}

func (p *PeerInfo) MakeMessage(m Message) {
	req := make([]byte, m.Length)
	switch m.ID {
	case Request:
		binary.BigEndian.PutUint32(req[0:], uint32(m.Length))
		req = append(req, byte(m.ID))
		req = append(req, m.Payload...)
		p.utp.BuildAndTransmitPacket(req)

	}

}
