package main

import (
	// "crypto/sha1"

	"fmt"

	// "io"
	"morpho/bencoding"
	"morpho/torrent"

	"os"
	// "strconv"
)

func main() {
	data, _ := os.ReadFile("another.torrent")
	source := string(data)
	bval, _ := bencoding.Decode(source)
	m := bval.(map[string]any)
	// fmt.Println("announce list ", m["announce-list"])

	meta_info, _ := torrent.LoadTorrent(bval)
	// fmt.Println(meta_info.AnnounceList[0])

	// conn, err := net.Dial("udp", "opentor.net:696913")
	// defer conn.Close()
	// if err != nil {
	// 	fmt.Println(err)
	// }

	// fmt.Println(m["announce-list"])

	// buf := make([]byte, 2048)
	// // go func() {

	// // }()
	announceData := torrent.CreateAnnounceData(&meta_info, m)
	// body := announceData.ToBytes()
	// bod, _ := announceData.ToHttp(&meta_info, *announce)
	// fmt.Println(string(bod))

	// w, err := conn.Write(body)
	// if err != nil {
	// 	fmt.Println("Error sending data:", err)
	// 	return
	// }
	// fmt.Println(w)
	// // for {
	// n, err := conn.Read(buf)
	// if err != nil {
	// 	fmt.Println(err)
	// 	if err != net.ErrClosed {
	// 		fmt.Println("Error reading data:", err)
	// 	}
	// 	conn.Close()
	// 	// break
	// }
	// msg := buf[:n]
	// fmt.Println(string(msg))

	//info sha hash
	// info := info_dic["info"]
	// encoded_info := bencoding.Encode(info)
	// h := sha1.New()
	// io.WriteString(h, encoded_info)

	// }
	dataChannel := make(chan interface{})

	body := announceData.ManageAnnounceTracker(&meta_info, dataChannel)
	// torrent.FromHTTP(body)
	fmt.Println("--------this is main---------", body, "---------------------------------")
	go torrent.ManageResponceData(dataChannel)

	// for i, v := range m {
	// 	fmt.Println(i, v)
	// 	//
	// }

}

// 	tracker, _ := bencoding.Decode(string(body))

// 	tm := tracker.(map[string]any)
// 	s := tm["peers"]
// 	// fmt.Println([]byte(s.(string)))
// 	fmt.Println(s)
// }
