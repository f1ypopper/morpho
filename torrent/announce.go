package torrent

import (
	"bytes"
	"context"
	"time"

	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"morpho/bencoding"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"sync"
	// "time"
)

func (ad *AnnounceData) ToBytes() []byte {
	/*
	   	Offset  Size    Name    Value
	   0       64-bit integer  connection_id
	   8       32-bit integer  action          1 // announce
	   12      32-bit integer  transaction_id
	   16      20-byte string  info_hash
	   36      20-byte string  peer_id
	   56      64-bit integer  downloaded
	   64      64-bit integer  left
	   72      64-bit integer  uploaded
	   80      32-bit integer  event           0 // 0: none; 1: completed; 2: started; 3: stopped
	   84      32-bit integer  IP address      0 // default
	   88      32-bit integer  key
	   92      32-bit integer  num_want        -1 // default
	   96      16-bit integer  port
	   98
	*/

	endian := binary.LittleEndian
	buf := new(bytes.Buffer)
	binary.Write(buf, endian, ad.ConnectionID)
	binary.Write(buf, endian, ad.Action)
	binary.Write(buf, endian, ad.TransactionID)
	binary.Write(buf, endian, []byte(ad.InfoHash))
	binary.Write(buf, endian, []byte(ad.PeerID))
	binary.Write(buf, endian, ad.Downloaded)
	binary.Write(buf, endian, ad.Left)
	binary.Write(buf, endian, ad.Uploaded)
	var event uint32 = 0
	switch ad.Event {
	case "completed":
		event = 1
	case "started":
		event = 2
	case "stopped":
		event = 3
	}
	binary.Write(buf, endian, event)
	binary.Write(buf, endian, uint32(0)) // IP is optional
	binary.Write(buf, endian, uint32(0)) // we don't store a key
	binary.Write(buf, endian, int32(-1))
	binary.Write(buf, endian, uint16(6881))
	return buf.Bytes()
}

func (a *AnnounceData) ToHttp(m *MetaInfo, announceUrl url.URL) ([]byte, error) {
	client := &http.Client{
		Timeout: 3 * time.Second,
	}
	params := url.Values{
		"info_hash":  {string(a.InfoHash)},
		"peer_id":    {a.PeerID},
		"port":       {announceUrl.Port()},
		"uploaded":   {strconv.Itoa(int(a.Uploaded))},
		"downloaded": {strconv.Itoa(int(a.Downloaded))},
		"left":       {strconv.Itoa(int(m.Info.Files[0].Length))},
		"compact":    { /*strconv.FormatBool(a.Compact)*/ "1"},
		"event":      {a.Event},
	}
	host := announceUrl.Host
	scheme := announceUrl.Scheme
	path := announceUrl.Path
	fullUrl := scheme + "://" + host + path + "?" + params.Encode()

	// req, err := http.NewRequest(http.MethodGet, fullUrl, nil)
	// if err != nil {
	// 	fmt.Println("Error creating request:", err)
	// 	return nil, err
	// }
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fullUrl, nil)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return nil, err
	}

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return nil, err
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	if body == nil {
		fmt.Println("Error in the body")
	}
	return body, nil
}

func ManageAnnounceList(aList []interface{}) []url.URL {
	var list []url.URL

	for _, v := range aList {
		if firstURL, ok := v.([]interface{}); ok {
			announce, _ := url.Parse(firstURL[0].(string))
			if announce.Scheme != "udp" {
				list = append(list, *announce)

			}

		} else {
			fmt.Println("Unexpected type in tracker list")
		}

	}
	return list
}

func (aData *AnnounceData) ManageAnnounceTracker(m *MetaInfo, peerList *[]Peer) {

	var wg sync.WaitGroup
	var mu sync.Mutex

	// for {
	for _, v := range m.AnnounceList {
		wg.Add(1)
		go func(aUrl *MetaInfo) ([]byte, error) {
			defer wg.Done()

			// _, cancel := context.WithTimeout(context.Background(), time.Second*3)
			// defer cancel()
			body, err := aData.ToHttp(m, v)
			if err != nil {
				return nil, err
			}
			t, _ := bencoding.Decode(string(body))
			if tracker, ok := t.(map[string]interface{}); ok {
				// fmt.Println("TRACKER ", tracker)
				if err, ok := tracker["failure reason"]; ok {
					return nil, errors.New(err.(string))
				}
				respData := FromHTTP(tracker)
				mu.Lock()

				*peerList = append(*peerList, respData.Peer...)
				mu.Unlock()
				// time.Sleep(time.Duration(respData.Interval) * time.Second)
			}
			return body, nil

		}(m)

		// }
		wg.Wait()
		// fmt.Println("completed reading trackers")
	}
}

func FromHTTP(tm map[string]interface{}) *ResponseData {
	var p []Peer
	switch val := tm["peers"].(type) {
	case string:
		bytesSlice := []byte(val)
		for offset := 0; offset < len(bytesSlice); offset += 6 {
			b := bytesSlice[offset : offset+6]
			if len(bytesSlice) > 0 {
				q := Peer{
					IP:   net.IP(b[:4]),
					Port: uint16(binary.BigEndian.Uint16(b[4:])),
				}
				p = append(p, q)

			}
		}
	case []interface{}:
		for _, v := range val { // [interface{}, intrface{}, interface{}]
			if peer, ok := v.(map[string]interface{}); ok {
				q := Peer{
					IP:     net.ParseIP(peer["ip"].(string)),
					PeerID: peer["peer id"].(string),
					Port:   uint16(peer["port"].(int)),
				}
				p = append(p, q)

			}
		}

	}
	returnData := &ResponseData{
		Complete:    uint(tm["complete"].(int)),
		Incomplete:  uint(tm["incomplete"].(int)),
		Interval:    uint(tm["interval"].(int)),
		MinInterval: uint(tm["min interval"].(int)),
		Peer:        p,
	}
	return returnData
}
