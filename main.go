package main

import (
	"crypto/sha1"
	"fmt"
	"io"
	"morpho/bencoding"
	"morpho/torrent"
	"net/http"
	"net/url"
	"os"
	"strconv"
)

func main() {
	data, _ := os.ReadFile("test.torrent")
	source := string(data)
	bval, _ := bencoding.Decode(source)
	m := bval.(map[string]any)
	info_dict := m["info"]
	encoded_info := bencoding.Encode(info_dict)
	h := sha1.New()
	io.WriteString(h, encoded_info)
	fmt.Printf("%x\n", string(h.Sum(nil)))
	fmt.Println(url.QueryEscape(string(h.Sum(nil))))
	meta_info, _ := torrent.LoadTorrent(bval)
	announce, _ := url.Parse(meta_info.Announce)
	q := announce.Query()
	q.Set("info_hash", string(h.Sum(nil)))
	q.Set("peer_id", "AAAAAAAAAAAAAAAAAAAA")
	q.Set("port", "6881")
	q.Set("uploaded", "0")
	q.Set("downloaded", "0")
	q.Set("left", strconv.Itoa(int(meta_info.Info.Files[0].Length)))
	q.Set("compact", "0")
	q.Set("event", "started")
	announce.RawQuery = q.Encode()
	res, err := http.Get(announce.String())
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()
	body, _ := io.ReadAll(res.Body)
	tracker, _ := bencoding.Decode(string(body))
	tm := tracker.(map[string]any)
	s := tm["peers"]
	fmt.Println([]byte(s.(string)))
}
