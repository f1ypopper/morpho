package torrent

import (
	"bytes"
	"crypto/sha1"
	"encoding/binary"
	"fmt"
	"io"
	"morpho/bencoding"
	"net/http"
	"net/url"
	"strconv"
	"sync"
)

func CreateAnnounceData(metainfo *MetaInfo, info_dic map[string]any) AnnounceData {

	info := info_dic["info"]
	encoded_info := bencoding.Encode(info)
	h := sha1.New()
	io.WriteString(h, encoded_info)
	// urlInfo, _ := url.Parse(metainfo.AnnounceURL)
	// portValue, _ := strconv.Atoi(urlInfo.Port())
	returnData := AnnounceData{
		Left:          uint64(metainfo.Info.Files[0].Length),
		PeerID:        "AAAAAAAAAAAAAAAAAAAA",
		Port:          80,
		Uploaded:      0,
		Downloaded:    0,
		Compact:       false,
		Event:         "started",
		ConnectionID:  0,
		Action:        1,
		TransactionID: 0,
		InfoHash:      h.Sum(nil),
	}
	return returnData

}

func LoadTorrent(bval any) (MetaInfo, error) {
	m, _ := bval.(map[string]any)
	metainfo := MetaInfo{}
	metainfo.AnnounceURL = m["announce"].(string)
	metainfo.AnnounceList = ManageAnnounceList(m["announce-list"].([]interface{}))
	metainfo.Info, _ = loadInfo(m["info"])
	return metainfo, nil
}

func loadInfo(bval any) (TorrentInfo, error) {
	m, _ := bval.(map[string]any)
	tinfo := TorrentInfo{}
	tinfo.PieceLength = uint(m["piece length"].(int))
	tinfo.Pieces = m["pieces"].(string)
	files := []File{}
	if files_dict, ok := m["files"]; ok {
		//multi-file mode
		for _, v := range files_dict.([]any) {
			file_dict, _ := v.(map[string]any)
			file := File{}
			file.Length = uint(file_dict["length"].(int))
			file.Path = []any{m["name"].(string)}
			files = append(files, file)
		}

	} else {
		file := File{}
		file.Length = uint(m["length"].(int))
		file.Path = []any{m["name"].(string)}
		files = append(files, file)
	}
	tinfo.Files = files
	return tinfo, nil
}

func (rd *ResponseData) PeerHandshake(info_dic map[string]any) []byte {
	/*
		<pstrlen><pstr><reserved><info_hash><peer_id>
			It is (49+len(pstr)) bytes long.

		pstrlen: 	string length of <pstr>, as a single raw byte
		pstr: 		string identifier of the protocol
		reserved: 	eight (8) reserved bytes. All current implementations use all zeroes.
		info_hash: 	20-byte SHA1 hash of the info key in the metainfo file.
		peer_id: 	20-byte string used as a unique ID for the client
	*/

	// info hash
	info := info_dic["info"]
	encoded_info := bencoding.Encode(info)
	h := sha1.New()
	io.WriteString(h, encoded_info)

	endian := binary.BigEndian
	buf := new(bytes.Buffer)
	binary.Write(buf, endian, 19)
	binary.Write(buf, endian, "BitTorrent protocol")
	binary.Write(buf, endian, int8(0))
	binary.Write(buf, endian, []byte(h.Sum(nil)))

	for _, peer := range rd.Peers {
		binary.Write(buf, endian, []byte(peer.PeerID))
		return buf.Bytes()
	}
	return nil
}

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
	// fmt.Println(len(buf.Bytes()))
	return buf.Bytes()
}

func (a *AnnounceData) ToHttp(m *MetaInfo, announceUrl url.URL) ([]byte, error) {
	client := &http.Client{}
	params := url.Values{
		"info_hash":  {string(a.InfoHash)},
		"peer_id":    {a.PeerID},
		"port":       {announceUrl.Port()},
		"uploaded":   {strconv.Itoa(int(a.Uploaded))},
		"downloaded": {strconv.Itoa(int(a.Downloaded))},
		"left":       {strconv.Itoa(int(m.Info.Files[0].Length))},
		"compact":    { /*strconv.FormatBool(a.Compact)*/ "0"},
		"event":      {a.Event},
	}
	// fullUrl := m.AnnounceURL + "?" + params.Encode()

	host := announceUrl.Host
	scheme := announceUrl.Scheme
	path := announceUrl.Path
	fullUrl := scheme + "://" + host + path + "?" + params.Encode()

	req, err := http.NewRequest(http.MethodGet, fullUrl, nil)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return nil, err
	}

	resp, err := client.Do(req)
	if err != nil {
		// fmt.Println("Error sending request:", err)
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
		// fmt.Printf("%T", v)
		if firstURL, ok := v.([]interface{}); ok {
			announce, _ := url.Parse(firstURL[0].(string))
			if announce.Scheme != "udp" {
				list = append(list, *announce)

			}

		} else {
			fmt.Println("Unexpected type in tracker list")
		}

	}
	fmt.Println("This is the length of the list ", list)
	return list

	//
}

// func (aData *AnnounceData) ManageAnnounceTracker(m *MetaInfo, dataChannel chan<- interface{} []byte {
func (aData *AnnounceData) ManageAnnounceTracker(m *MetaInfo, dataChannel chan<- interface{}) []byte {
	var wg sync.WaitGroup

	for i, v := range m.AnnounceList {
		wg.Add(1)
		go func(aUrl *MetaInfo) ([]byte, error) {
			defer wg.Done()
			// fmt.Println(v)
			body, err := aData.ToHttp(m, v)
			if err != nil {
				copy(m.AnnounceList[i:], m.AnnounceList[i+1:])

				m.AnnounceList = m.AnnounceList[:len(m.AnnounceList)-1]
				fmt.Println("this is announce list ", m.AnnounceList)
				fmt.Println("The error is ", err)
				return nil, err

			}
			// fmt.Println(string(body))
			tracker, _ := bencoding.Decode(string(body))
			if tracker != nil {
				// fmt.Printf("____bencoded tracker is %T : %v ____\n", tracker, tracker)
				if _, ok := tracker.(map[string]interface{}); ok {
					fmt.Println(ok)

					FromHTTP(tracker.(map[string]interface{}))
				}
				dataChannel <- tracker
			}
			return body, nil

			// fmt.Println("pinging ", aUrl.Host())

		}(m)

	}
	wg.Wait()
	return nil
}

func FromHTTP(tm map[string]interface{}) ResponseData {

	var p []Peers
	if _, ok := tm["peers"]; ok {
		for _, v := range tm["peers"].([]interface{}) { // [interface{}, intrface{}, interface{}]
			if peer, ok := v.(map[string]interface{}); ok {

				q := Peers{
					IP:     peer["ip"].(string),
					PeerID: peer["peer id"].(string),
					Port:   uint16(peer["port"].(int)),
				}
				p = append(p, q)

			}
		}
	} else {
		fmt.Printf("peers - interface {} is nil.\n ")
	}
	fmt.Printf("%T is type of complete    \n", tm["complete"])

	returnData := ResponseData{
		Complete:    uint(tm["complete"].(int)),
		Incomplete:  uint(tm["incomplete"].(int)),
		Interval:    uint(tm["interval"].(int)),
		MinInterval: uint(tm["min interval"].(int)),
		Peers:       p,
	}
	fmt.Printf("this is the return  data - %T : %v  \n", returnData, returnData)

	return returnData
}

func ManageResponceData(dataChannel <-chan interface{}) {

	responseData := <-dataChannel

	fmt.Println("responce data throucg channel", responseData)
}
