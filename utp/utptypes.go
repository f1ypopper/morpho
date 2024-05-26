package utp

import "net"

type ConState int

const (
	CS_UNINITIALIZED ConState = 0
	CS_SYN_SENT
	CS_CONNECTED
)

type UTPConnection struct {
	BaseConn     net.Conn
	ip           string
	state        ConState
	seqNr        uint16
	ackNr        uint16
	conn_id_recv uint16
	conn_id_send uint16
}
type Packet struct {
	ptype                             uint16
	connection_id                     uint16
	timestamp_microseconds            uint32
	timestamp_difference_microseconds uint32
	wnd_size                          uint32
	seq_nr                            uint16
	ack_nr                            uint16
	Payload                           []byte
}
