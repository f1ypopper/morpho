package utp

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

type PacketType uint32

const (
	ST_DATA  PacketType = 0
	ST_FIN              = 1
	ST_STATE            = 2
	ST_RESET            = 3
	ST_SYN              = 4
)

type Packet struct {
	ptype                             PacketType
	connection_id                     uint16
	timestamp_microseconds            uint32
	timestamp_difference_microseconds uint32
	wnd_size                          uint32 //their advertised window size
	seq_nr                            uint16
	ack_nr                            uint16
	payload                           []byte
}

func (packet *Packet) serialize() []byte {
	//pack into a single uint32
	var first uint32 = 0
	first |= uint32(packet.ptype << 28)
	first |= uint32(1) << 24
	first |= uint32(0) << 20
	first |= uint32(packet.connection_id)
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.BigEndian, first)
	binary.Write(buf, binary.BigEndian, packet.timestamp_microseconds)
	binary.Write(buf, binary.BigEndian, packet.timestamp_difference_microseconds)
	binary.Write(buf, binary.BigEndian, packet.wnd_size)
	var seq_ack_nr uint32 = 0
	seq_ack_nr |= uint32(packet.seq_nr) << 16
	seq_ack_nr |= uint32(packet.ack_nr)
	binary.Write(buf, binary.BigEndian, seq_ack_nr)
	if packet.payload != nil {
		binary.Write(buf, binary.BigEndian, packet.payload)
	}
	return buf.Bytes()
}

func (packet *Packet) deserialize(buf []byte) {
	first := binary.BigEndian.Uint32(buf[:4])
	packet.ptype = PacketType(first >> 28)
	packet.connection_id = binary.BigEndian.Uint16(buf[2:4])
	packet.timestamp_microseconds = binary.BigEndian.Uint32(buf[4:8])
	packet.timestamp_difference_microseconds = binary.BigEndian.Uint32(buf[8:12])
	packet.wnd_size = binary.BigEndian.Uint32(buf[12:16])
	var seq_ack_nr = binary.BigEndian.Uint32(buf[16:20])
	packet.seq_nr = uint16(seq_ack_nr >> 16)
	packet.ack_nr = uint16(seq_ack_nr & 0xFFFF)
	packet.payload = buf[20:]
}

func (p *Packet) len() uint32 {
	return 20 + uint32(len(p.payload))
}

func (p *Packet) to_str() string {
	var ptype string
	switch p.ptype {
	case ST_DATA:
		ptype = "ST_DATA"
	case ST_STATE:
		ptype = "ST_STATE"
	case ST_FIN:
		ptype = "ST_FIN"
	case ST_RESET:
		ptype = "ST_RESET"
	case ST_SYN:
		ptype = "ST_SYN"
	}
	return fmt.Sprintf("TYPE: %s ACK_NR: %d SEQ_NR: %d PAYLOAD: %s", ptype, p.ack_nr, p.seq_nr, p.payload)
}
