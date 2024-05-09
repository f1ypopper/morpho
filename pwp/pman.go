package pwp

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"morpho/torrent"
	"morpho/utp"
	"strconv"
	"sync"
)

type PeerInfo struct {
	interested bool
}

func handlePeer() {
}

type PieceInfo struct {
	state string //COMPLETE, INCOMPLETE, PROGRESS
}

type PeerConn struct {
	//conn net.Conn
	conn utp.UTPConnection
}

func newPeerConn(ip string, port uint16) (PeerConn, error) {
	//conn, err := net.Dial("tcp", ip+":"+strconv.Itoa(int(port)))
	conn, err := utp.Dial(ip + ":" + strconv.Itoa(int(port)))
	if err != nil {
		return PeerConn{}, err
	}
	return PeerConn{conn}, err
}

func (pconn *PeerConn) Handshake(peer_id []byte, info_hash []byte) ([]byte, error) {
	endian := binary.BigEndian
	wbuf := new(bytes.Buffer)

	var pstrlen uint8 = 19
	pstr := []byte("BitTorrent protocol")
	var reserved uint64 = 0
	binary.Write(wbuf, endian, pstrlen)
	binary.Write(wbuf, endian, pstr)
	binary.Write(wbuf, endian, reserved)
	binary.Write(wbuf, endian, info_hash)
	binary.Write(wbuf, endian, peer_id)

	if wbuf.Len() != int(pstrlen)+49 {
		return nil, errors.New("wrong length of wbuffer")
	}
	bytes_written, err := pconn.conn.Write(wbuf.Bytes())
	if err != nil {
		return nil, err
	}
	if bytes_written != int(pstrlen)+49 {
		return nil, fmt.Errorf("incomplete bytes written, required: %d written: %d", int(pstrlen)+49, bytes_written)
	}

	var rbuf = make([]byte, 49+int(pstrlen))
	_, err = pconn.conn.Read(rbuf)
	if err != nil {
		return nil, err
	}
	return rbuf, nil
}

func StartPeerManager(peermap *map[string]uint16, adata *torrent.AnnounceData) {
	//fmt.Println(peermap)
	var wg sync.WaitGroup
	for ip, port := range *peermap {
		wg.Add(1)
		go func() {
			defer wg.Done()
			conn, err := newPeerConn(ip, port)
			if err != nil {
				fmt.Println(err)
				return
			}
			fmt.Printf("%s:%d CONNECTED\n", ip, port)
			rbuf, err := conn.Handshake([]byte(adata.PeerID), adata.InfoHash)
			if err != nil {
				fmt.Println("ERROR IN PWP HANDSHAKE: ", err)
			}
			fmt.Printf("%s:%d HANDSHAKE: %s\n", ip, port, rbuf)
			//conn.Handshake([]byte(adata.PeerID), adata.InfoHash)
		}()
	}
	wg.Wait()
}
