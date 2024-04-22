package utp

import (
	"errors"
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
const MAX_PAYLOAD_LEN = 1004

type ConnState int

const (
	CS_UNINITIALIZED ConnState = 0
	CS_SYN_SENT
	CS_CONNECTED
)

type UTPConnection struct {
	baseconn      net.Conn //maybe switch to net.UDPConn
	id_send       uint16   // send connection_id
	id_recv       uint16   // recv connection_id
	seq_nr        uint16
	ack_nr        uint16 //last acked message
	cur_window    uint   //bytes in flight (sent but not acked)
	max_wind_size uint32
	their_wnd     uint32
	state         ConnState
	readbuf       []byte //bytes recieved
}

func Dial(address string) (UTPConnection, error) {
	baseconn, err := net.Dial("udp", address)
	if err != nil {
		return UTPConnection{}, nil
	}
	//TODO handshake
	conn := UTPConnection{baseconn: baseconn, state: CS_UNINITIALIZED, readbuf: make([]byte, 0, 1000000)}
	if err := conn.syn(); err != nil {
		return UTPConnection{}, err
	}
	return conn, nil
}

func (conn *UTPConnection) Write(b []byte) (int, error) {
	p := Packet{}
	p.ptype = ST_DATA
	p.connection_id = conn.id_send
	p.seq_nr = conn.seq_nr
	p.ack_nr = conn.ack_nr
	p.timestamp_microseconds = uint32(time.Now().UnixMicro())
	p.timestamp_difference_microseconds = 0
	p.wnd_size = 1048576
	p.payload = b
	conn.send_packet(&p)
	conn.seq_nr += 1
	return len(b), nil
}

func (conn *UTPConnection) Read(b []byte) (int, error) {
	for len(b) > len(conn.readbuf) {
		packet, err := conn.recv_packet()
		if err != nil {
			return 0, err
		}
		conn.process_packet(&packet)
	}
	return copy(b, conn.readbuf), nil
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
	packet.timestamp_microseconds = uint32(time.Now().UnixMicro())
	conn.send_packet(&packet)
	ack, err := conn.recv_packet()
	if err != nil {
		return err
	}
	conn.process_packet(&ack)
	if conn.state != CS_CONNECTED {
		return errors.New("failed syn-ack handshake")
	}
	conn.seq_nr += 1
	return nil
}

func (conn *UTPConnection) process_packet(packet *Packet) error {
	switch packet.ptype {
	case ST_STATE:
		{
			if conn.state == CS_SYN_SENT {
				conn.state = CS_CONNECTED
				conn.ack_nr = packet.seq_nr - 1
			}
			if len(packet.payload) != 0 {
				conn.ack_nr = packet.seq_nr - 1
			}
		}
	case ST_DATA:
		{
			//if seq_nr is not conn.ack_nr+1 drop the packet (since we haven't created a priority queue for the ordering of packets)
			conn.readbuf = append(conn.readbuf, packet.payload...)
			conn.ack(packet.seq_nr)
			conn.ack_nr = packet.seq_nr
		}
	}
	return nil
}

func (conn *UTPConnection) ack(last_ack uint16) {
	p := Packet{}
	p.ptype = ST_STATE
	p.connection_id = conn.id_send
	p.seq_nr = conn.seq_nr
	p.ack_nr = last_ack
	p.timestamp_microseconds = uint32(time.Now().UnixMicro())
	p.timestamp_difference_microseconds = 0
	p.wnd_size = 1048576
	conn.baseconn.Write(p.serialize())
}

func (conn *UTPConnection) send_packet(packet *Packet) {
	if uint32(conn.cur_window)+packet.len() > min(conn.their_wnd, uint32(conn.cur_window)) {
		//wait for ack messages i.e wait for endpoint to process packets from its recive buffer
	}
	buf := packet.serialize()
	conn.baseconn.Write(buf)
	conn.cur_window += uint(len(buf))
}

func (conn *UTPConnection) recv_packet() (Packet, error) {
	buf := make([]byte, MAX_PAYLOAD_LEN)
	conn.baseconn.Read(buf)
	packet := Packet{}
	packet.deserialize(buf)
	return packet, nil
}
