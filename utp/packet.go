package utp

import (
	"context"
	"encoding/binary"
	"fmt"
	"net"
	"strings"
	"time"
)

func (c *UTPConnection) BuildAndTransmitPacket(payloadData []byte) Packet {
	var p Packet
	p.ptype = 0
	p.seq_nr = c.seqNr
	p.ack_nr = c.ackNr
	p.connection_id = c.conn_id_send
	p.timestamp_microseconds = uint32(time.Now().UnixMicro())
	p.timestamp_difference_microseconds = 0
	p.wnd_size = 1048576
	p.Payload = payloadData
	// send packet
	p.SendPacket(c.BaseConn)
	return p
}

// Send the packet by writing to the connection.
func (packet *Packet) SendPacket(conn net.Conn) {
	buf := packet.serialize()
	conn.Write(buf)

}

// handshake with retrieving bitfield
func (c *UTPConnection) HandshakePacket(ctx context.Context, payloadData []byte) ([]byte, error) {
	c.BuildAndTransmitPacket(payloadData)
	timeout := time.Now().Add(10 * time.Second)

	for time.Now().Before(timeout) {
		select {
		case <-ctx.Done():
			fmt.Println("Exiting due to context timeout")
			return nil, ctx.Err()
		default:
			res_buf := make([]byte, 1024)
			// c.BaseConn.SetReadDeadline(time.Now().Add(time.Second * 4))
			n, err := c.BaseConn.Read(res_buf)
			if err != nil {
				return nil, err
			}
			packet := deserialize(res_buf[:n])
			if strings.Contains(string(packet.Payload), "BitTorrent protocol") {
				peers = append(peers, c.ip)
				len, msg_id := binary.BigEndian.Uint32(res_buf[88:92]), res_buf[92]
				if msg_id == 5 {
					bitfield := res_buf[93 : 92+len]
					return bitfield, nil
				}

			}
			if time.Now().After(timeout) {
				fmt.Println("exiting timeout")
				return nil, nil
			}
		}

	}
	return nil, nil
}
