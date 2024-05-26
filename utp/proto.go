package utp

import (
	"context"
	"fmt"
	"math/rand"
	"net"
	"time"
)

var peers []string

func (c *UTPConnection) ProcessReceivedPacket(res []byte) Packet {

	if len(res) > 0 {
		recvPacket := deserialize(res)
		if recvPacket.ack_nr == c.seqNr && recvPacket.ptype == 2 {
			c.state = CS_CONNECTED
			c.ackNr = recvPacket.seq_nr - 1
			c.seqNr += 1
		}
		return recvPacket
	}
	return Packet{}
}

// INITIALIZE CONNECTION
func InitConnection(ctx context.Context, ip string, timeout time.Duration) (net.Conn, error) {

	conn, err := net.DialTimeout("udp", ip, timeout*time.Second)
	if err != nil {
		fmt.Println("net dial error ")
		return nil, err
	}
	select {
	case <-ctx.Done():
		fmt.Println("Exiting due to context timeout")
		return nil, err
	default:
		// Continue the loop

	}
	return conn, nil

}

func (c *UTPConnection) Syn() error {
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
	p.SendPacket(c.BaseConn)
	res_buf := make([]byte, 1024)
	// c.BaseConn.SetReadDeadline(time.Now().Add(5 * time.Second))
	n, err := c.BaseConn.Read(res_buf)
	if err != nil {
		return err

	}

	if n > 0 {
		c.ProcessReceivedPacket(res_buf[:n])

	}

	return nil

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
	p.SendPacket(c.BaseConn)
	res_buf := make([]byte, 1024)
	c.BaseConn.Read(res_buf)

}

// optimize this
func NewPeer(ipAddr string, handshake []byte) ([]byte, UTPConnection, error) {
	var c UTPConnection
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	timeout := time.Now().Add(10 * time.Second)

	defer cancel()
	c.BaseConn, _ = InitConnection(ctx, ipAddr, 10)
	if c.BaseConn == nil {
		fmt.Println("base conn is nil")
		return nil, UTPConnection{}, fmt.Errorf("base conn is nil")
	}
	c.ip = ipAddr
	err := c.Syn()
	if err != nil {
		return nil, UTPConnection{}, err
	}
	bitfield, err := c.HandshakePacket(ctx, handshake)
	if err != nil {
		if err == context.DeadlineExceeded {
			fmt.Println("Handshake packet timed out")
			return nil, UTPConnection{}, err

		} else {
			fmt.Println("Error in handshake packet:", err)
			return nil, UTPConnection{}, err

		}
	}
	if time.Now().After(timeout) {
		fmt.Println("exiting timeout")
		return nil, UTPConnection{}, err
	}

	return bitfield, c, nil

}
