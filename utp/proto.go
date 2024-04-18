package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"math/rand"
	"net"
	"time"
)

/*
	HEADER (20 bytes)
0       4       8               16              24              32
+-------+-------+---------------+---------------+---------------+
| type  | ver   | extension     | connection_id                 |
+-------+-------+---------------+---------------+---------------+
| timestamp_microseconds                                        |
+---------------+---------------+---------------+---------------+
| timestamp_difference_microseconds                             |
+---------------+---------------+---------------+---------------+
| wnd_size                                                      |
+---------------+---------------+---------------+---------------+
| seq_nr                        | ack_nr                        |
+---------------+---------------+---------------+---------------+
*/

const HEADER_LEN = 20
const MAX_PACKET_LEN = 1480

type ConnState int

const (
	CS_UNINITIALIZED ConnState = 0
	CS_SYN_SENT
)

type UTPConnection struct {
	baseconn      net.Conn //maybe switch to net.UDPConn
	id_send       uint16   // send connection_id
	id_recv       uint16   // recv connection_id
	seq_nr        uint16
	ack_nr        uint16
	cur_window    uint16
	max_wind_size uint32
	their_wnd     uint32
	state         ConnState
}

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

func (p *Packet) serialize() []byte {
	panic("todo")
}

func (p *Packet) deserialize([]byte) {
	panic("todo")
}

func (p *Packet) len() uint32 {
	return 20
}

func Dial(address string) (UTPConnection, error) {
	baseconn, err := net.Dial("udp", address)
	if err != nil {
		return UTPConnection{}, nil
	}
	//TODO handshake
	conn := UTPConnection{baseconn: baseconn, state: CS_UNINITIALIZED}
	if err := conn.syn(); err != nil {
		return UTPConnection{}, err
	}
	return conn, nil
}

func (conn *UTPConnection) Write(b []byte) (int, error) {
	return 0, nil
}

func (conn *UTPConnection) Read(b []byte) (int, error) {
	return 0, nil
}

func (conn *UTPConnection) syn() error {
	//send the syn packet and wait for the ack
	conn.seq_nr = 1
	conn.id_recv = uint16(rand.Int())
	conn.id_send = conn.id_recv + 1
	packet := Packet{}
	packet.ptype = ST_SYN
	packet.connection_id = conn.id_recv
	packet.ack_nr = 0
	packet.seq_nr = conn.seq_nr
	packet.timestamp_difference_microseconds = 0
	conn.send_packet(&packet)
	return nil
}

func (conn *UTPConnection) send_packet(packet *Packet) {
	if uint32(conn.cur_window)+packet.len() > min(conn.their_wnd, uint32(conn.cur_window)) {
		//wait for ack messages
	}
	//pack into a single uint32
	var first uint32 = 0
	first |= uint32(packet.ptype << 28)
	first |= uint32(1) << 24
	first |= uint32(0) << 20
	first |= uint32(packet.connection_id)
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.BigEndian, first)
	binary.Write(buf, binary.BigEndian, uint32(time.Now().UnixMicro()))
	binary.Write(buf, binary.BigEndian, uint32(0)) //timestamp_difference_microseconds
	binary.Write(buf, binary.BigEndian, packet.wnd_size)
	var seq_ack_nr uint32 = 0
	seq_ack_nr |= uint32(packet.seq_nr)
	seq_ack_nr |= uint32(packet.ack_nr)
	binary.Write(buf, binary.BigEndian, seq_ack_nr)
	conn.baseconn.Write(buf.Bytes())
	res_buffer := make([]byte, 20)
	conn.baseconn.Read(res_buffer)
	fmt.Printf("BYTES RECIEVED: %x", res_buffer)
}

func main() {
	_, err := Dial("localhost:1111")
	if err != nil {
		fmt.Printf("CONNECTION ERR: %e\n", err)
	}
}
