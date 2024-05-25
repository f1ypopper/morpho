package utp

// import (
// 	"encoding/binary"
// 	"fmt"
// 	"strings"
// )

// func (c *UTPConnection) HandshakePacket(payloadData []byte) {
// 	var p Packet
// 	c.MakePacket(payloadData)
// 	p.SendPacket(c.baseConn)

// 	for {
// 		res_buf := make([]byte, 250)
// 		n, err := c.baseConn.Read(res_buf)
// 		if err != nil {
// 			break
// 		}
// 		packet := deserialize(res_buf[:n])
// 		fmt.Println(packet)
// 		if strings.Contains(string(packet.payload), "BitTorrent protocol") {
// 			peers = append(peers, c.ip)
// 			len, msg_id := binary.BigEndian.Uint32(res_buf[88:92]), res_buf[92]
// 			if msg_id == 5 {
// 				bitfield := res_buf[93 : 92+len]
// 				fmt.Printf("LENGTH: %d MESSAGE ID: %d \n", len, msg_id)
// 				fmt.Println("BITFIELD:  ")
// 				for _, n := range bitfield {
// 					fmt.Printf("%08b ", n)

// 				}
// 				fmt.Println("rest is ", res_buf[92+len:])
// 			}

// 			break
// 		}

// 	}
// }
