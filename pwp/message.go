package pwp

import (
	"encoding/binary"
	"morpho/torrent"
)

// func request message
// func choke and unchoke as well handle it
// request: <len=0013><id=6><index><begin><length>

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

func (p *PeerInfo) RequestMessage(index int, metaInfo *torrent.MetaInfo) {

	var m Message
	peiceSize := metaInfo.Info.PieceLength
	// totalsize := metaInfo.Info.Files[0].Length
	MaxPeiceSize := 16384 // 16 KB
	begin := index * int(peiceSize)
	end := begin + int(peiceSize)
	for i := 0; i < end/MaxPeiceSize; i++ {
		req := make([]byte, 13)
		binary.BigEndian.PutUint32(req[0:], uint32(13))
		req = append(req, byte(m.ID))
		// craft payload message
		payload := make([]byte, 12)

		binary.BigEndian.PutUint32(payload[0:], uint32(index))
		binary.BigEndian.PutUint32(payload[4:], uint32(begin))
		binary.BigEndian.PutUint32(payload[8:], uint32(MaxPeiceSize))
		req = append(req, payload...)
		p.utp.BuildAndTransmitPacket(req)

		begin += MaxPeiceSize
		end += MaxPeiceSize

	}
}

// msg can be either choke unchoke interested notinterest
func (p *PeerInfo) PwpMessage(msg Id) {
	mes := make([]byte, 5)

	binary.BigEndian.PutUint32(mes[0:], uint32(1))
	mes = append(mes, byte(msg))
	p.utp.BuildAndTransmitPacket(mes)
}
