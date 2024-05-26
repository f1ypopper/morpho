package utp

import (
	"bytes"
	"encoding/binary"
)

func deserialize(resBuf []byte) Packet {
	var packet Packet
	first := binary.BigEndian.Uint32(resBuf[:4])
	packet.ptype = uint16(first >> 28)
	packet.connection_id = binary.BigEndian.Uint16(resBuf[2:4])
	packet.timestamp_microseconds = binary.BigEndian.Uint32(resBuf[4:8])
	packet.timestamp_difference_microseconds = binary.BigEndian.Uint32(resBuf[8:12])
	packet.wnd_size = binary.BigEndian.Uint32(resBuf[12:16])
	var seq_ack_nr = binary.BigEndian.Uint32(resBuf[16:20])
	packet.seq_nr = uint16(seq_ack_nr >> 16)
	packet.ack_nr = uint16(seq_ack_nr & 0xFFFF)
	if len(resBuf) > 20 {
		packet.Payload = resBuf[20:]

	}

	return packet

}

func (pack *Packet) serialize() []byte {
	// get the len of payload and add to make|
	buf := make([]byte, 20)
	buffer := new(bytes.Buffer)
	var first uint16 = 0
	first |= uint16(pack.ptype << 12)
	first |= uint16(1 << 8)
	first |= uint16(uint8(0) << 0)
	binary.BigEndian.PutUint16(buf[0:], uint16(first))
	binary.BigEndian.PutUint16(buf[2:], pack.connection_id)
	binary.BigEndian.PutUint32(buf[4:], pack.timestamp_microseconds)
	binary.BigEndian.PutUint32(buf[8:], pack.timestamp_difference_microseconds)
	binary.BigEndian.PutUint32(buf[12:], pack.wnd_size)
	binary.BigEndian.PutUint16(buf[16:], pack.seq_nr)
	binary.BigEndian.PutUint16(buf[18:], pack.ack_nr)
	buf = append(buf, pack.Payload...)
	if pack.Payload != nil {
		binary.Write(buffer, binary.BigEndian, pack.Payload)
	}
	return buf
}
