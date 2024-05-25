package utp

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"math/rand"
	"net"
	"strings"
	"sync"
	"time"
)

// we send syn
// they send ack
// we send data ack their seq nr
// they recive our packet and set seq nr to our seq nr
// TODO
type ConState int

const (
	CS_UNINITIALIZED ConState = 0
	CS_SYN_SENT
	CS_CONNECTED
)

type UTPConnection struct {
	baseConn     net.Conn
	ip           string
	state        ConState
	seqNr        uint16
	ackNr        uint16
	conn_id_recv uint16
	conn_id_send uint16
	rbuflock     *sync.RWMutex
	readbuf      *bytes.Buffer
}
type Packet struct {
	ptype                             uint16
	connection_id                     uint16
	timestamp_microseconds            uint32
	timestamp_difference_microseconds uint32
	wnd_size                          uint32
	seq_nr                            uint16
	ack_nr                            uint16
	payload                           []byte
}

var peers []string

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
		packet.payload = resBuf[20:]

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
	buf = append(buf, pack.payload...)
	if pack.payload != nil {
		binary.Write(buffer, binary.BigEndian, pack.payload)
	}
	return buf
}
func (c *UTPConnection) CheckAcked(res []byte) {
	// check if packet is acked
	if len(res) > 0 {
		recvPacket := deserialize(res)
		if strings.Contains(string(recvPacket.payload), "BitTorrent protocol") {
			peers = append(peers, c.ip)

		}
		if recvPacket.ack_nr == c.seqNr && recvPacket.ptype == 2 {
			c.state = CS_CONNECTED
			c.ackNr = recvPacket.seq_nr - 1
			c.seqNr += 1
		}
	}
}

// INITIALIZE CONNECTION
func InitConnection(ip string, timeout time.Duration) net.Conn {

	conn, err := net.DialTimeout("udp", ip, timeout*time.Second)
	if err != nil {
		fmt.Println("net dial error ")
	}
	return conn

}

func (c *UTPConnection) Syn() {
	c.seqNr = 1
	c.ackNr = 0
	c.state = CS_SYN_SENT
	c.conn_id_recv = uint16(rand.Int())
	c.conn_id_send = c.conn_id_recv + 1
	// craft the packet.
	var p Packet
	p.ptype = 4
	p.seq_nr = c.seqNr
	p.ack_nr = c.ackNr
	p.connection_id = c.conn_id_recv
	p.timestamp_microseconds = uint32(time.Now().UnixMicro())
	p.timestamp_difference_microseconds = 0
	p.wnd_size = 0
	p.SendPacket(c.baseConn)
	res_buf := make([]byte, 1024)
	n, _ := c.baseConn.Read(res_buf)
	time.Sleep(time.Second)

	if n > 0 {
		c.CheckAcked(res_buf[:n])

	}

}
func (c *UTPConnection) Ack() {
	var p Packet
	p.ptype = 0
	p.seq_nr = c.seqNr
	p.ack_nr = 0
	p.connection_id = c.conn_id_send
	p.timestamp_microseconds = uint32(time.Now().UnixMicro())
	p.timestamp_difference_microseconds = 0
	p.wnd_size = 1048576
	p.SendPacket(c.baseConn)
	res_buf := make([]byte, 1024)
	c.baseConn.Read(res_buf)

}

// handshake with retrieving bitfield
func (c *UTPConnection) HandshakePacket(payloadData []byte) {
	var p Packet
	p.ptype = 0
	p.seq_nr = c.seqNr
	p.ack_nr = c.ackNr
	p.connection_id = c.conn_id_send
	p.timestamp_microseconds = uint32(time.Now().UnixMicro())
	p.timestamp_difference_microseconds = 0
	p.wnd_size = 1048576
	p.payload = payloadData
	p.SendPacket(c.baseConn)

	for {
		res_buf := make([]byte, 250)
		n, err := c.baseConn.Read(res_buf)
		if err != nil {
			break
		}
		packet := deserialize(res_buf[:n])
		if strings.Contains(string(packet.payload), "BitTorrent protocol") {
			peers = append(peers, c.ip)
			len, msg_id := binary.BigEndian.Uint32(res_buf[88:92]), res_buf[92]
			if msg_id == 5 {
				bitfield := res_buf[93 : 92+len]
				fmt.Printf("LENGTH: %d MESSAGE ID: %d \n", len, msg_id)
				fmt.Println("BITFIELD:  ")
				for _, n := range bitfield {
					fmt.Printf("%08b ", n)

				}
				fmt.Println("rest is ", res_buf[92+len:])
			}

			break
		}

	}
}

type messageID uint8

const (
	MsgChoke         messageID = 0
	MsgUnchoke       messageID = 1
	MsgInterested    messageID = 2
	MsgNotInterested messageID = 3
	MsgHave          messageID = 4
	MsgBitfield      messageID = 5
	MsgRequest       messageID = 6
	MsgPiece         messageID = 7
	MsgCancel        messageID = 8
)

// Message stores ID and payload of a message
type Message struct {
	ID      messageID
	Payload []byte
}

func (m *Message) message() []byte {
	if m == nil {
		return make([]byte, 235)
	}
	length := uint32(len(m.Payload) + 1) // +1 for id
	buf := make([]byte, 4+length)
	binary.BigEndian.PutUint32(buf[0:4], length)
	buf[4] = byte(m.ID)
	copy(buf[5:], m.Payload)
	fmt.Println(buf)
	return buf

}

// Send the packet by writing to the connection.
func (packet *Packet) SendPacket(conn net.Conn) {
	buf := packet.serialize()
	conn.Write(buf)

}
func Send(ipAddr string, handshake []byte) {
	var c UTPConnection
	c.baseConn = InitConnection(ipAddr, 1000)
	c.ip = ipAddr
	c.Syn()
	c.HandshakePacket(handshake)
}
