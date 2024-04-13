package main

import (
	// "crypto/sha1"
	"fmt"
	// "io"
	"morpho/bencoding"
	"morpho/torrent"
	"net"
	"net/url"
	"os"
	// "strconv"
)

func main() {
	data, _ := os.ReadFile("sample.torrent")
	source := string(data)
	bval, _ := bencoding.Decode(source)
	m := bval.(map[string]any)

	meta_info, _ := torrent.LoadTorrent(bval)
	announce, _ := url.Parse(meta_info.AnnounceURL)
	fmt.Println(announce.Port)
	conn, err := net.Dial("udp", "tracker.openbittorrent.com:8013")
	defer conn.Close()
	if err != nil {
		fmt.Println(err)
	}

	buf := make([]byte, 2048)
	go func() {
		for {
			n, err := conn.Read(buf)
			if err != nil {
				fmt.Println(err)
				if err != net.ErrClosed { // Handle other errors
					fmt.Println("Error reading data:", err)
				}
				conn.Close()
				break
			}
			msg := buf[:n]
			// fmt.Println(msg)
			fmt.Println(string(msg))

		}

	}()
	announceData := torrent.CreateAnnounceData(&meta_info, m)
	body := announceData.ToBytes()

	w, err := conn.Write(body)
	if err != nil {
		fmt.Println("Error sending data:", err)
		return
	}
	fmt.Println(w)

	tracker, _ := bencoding.Decode(string(body))
	fmt.Println(tracker)
	// tm := tracker.(map[string]any)
	// s := tm["peers"]
	// fmt.Println([]byte(s.(string)))
	// fmt.Println(string(s.(string)))
}
