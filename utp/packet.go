package utp

import (
	"context"
	"encoding/binary"
	"fmt"
	"math/rand"
	"strings"
	"time"
)

// import (
//
//	"encoding/binary"
//	"fmt"
//	"strings"
//
// )

func SetConnection() UTPConnection {
	var c UTPConnection
	c.seqNr = 1
	c.ackNr = 0
	c.state = CS_CONNECTED
	c.conn_id_recv = uint16(rand.Int())
	c.conn_id_send = c.conn_id_recv + 1
	return c

}

func (c *UTPConnection) BuildAndTransmitPacket(payloadData []byte) Packet {
	// c := SetConnection()
	var p Packet
	p.ptype = 0
	p.seq_nr = c.seqNr
	p.ack_nr = c.ackNr
	p.connection_id = c.conn_id_send
	p.timestamp_microseconds = uint32(time.Now().UnixMicro())
	p.timestamp_difference_microseconds = 0
	p.wnd_size = 1048576
	p.payload = payloadData
	// send packet
	p.SendPacket(c.baseConn)
	return p
}

// handshake with retrieving bitfield
func (c *UTPConnection) HandshakePacket(ctx context.Context, payloadData []byte) ([]byte, error) {
	c.BuildAndTransmitPacket(payloadData)
	timeout := time.Now().Add(10 * time.Second)
	done := make(chan bool)

	go func() {
		time.Sleep(time.Second * 10)
		done <- true
	}()

	for time.Now().Before(timeout) {
		select {
		case <-ctx.Done():
			fmt.Println("Exiting due to context timeout")
			return nil, ctx.Err()
		default:
			res_buf := make([]byte, 1024)
			c.baseConn.SetReadDeadline(time.Now().Add(time.Second * 1))
			n, err := c.baseConn.Read(res_buf)
			if err != nil {
				return nil, err
			}
			packet := deserialize(res_buf[:n])
			if strings.Contains(string(packet.payload), "BitTorrent protocol") {
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
