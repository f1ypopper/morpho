package torrent

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"hash"
	"net/url"
	"strconv"
)

func LoadTorrent(bval any) (MetaInfo, error) {
	m, _ := bval.(map[string]any)
	metainfo := MetaInfo{}
	metainfo.AnnounceURL = m["announce"].(string)
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
			file.Path = file_dict["path"].(string)
			files = append(files, file)
		}

	} else {
		file := File{}
		file.Length = uint(m["length"].(int))
		file.Path = m["name"].(string)
		files = append(files, file)
	}
	tinfo.Files = files
	return tinfo, nil
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
	fmt.Println(len(buf.Bytes()))
	return buf.Bytes()
}

func (a *AnnounceData) ToHttp(m *MetaInfo, hash hash.Hash) string {
	params := url.Values{
		"info_hash":  {string(hash.Sum(nil))},
		"peer_id":    {a.PeerID},
		"port":       {strconv.Itoa(int(a.Port))},
		"uploaded":   {strconv.Itoa(int(a.Uploaded))},
		"downloaded": {strconv.Itoa(int(a.Downloaded))},
		"left":       {strconv.Itoa(int(m.Info.Files[0].Length))},
		"compact":    {strconv.FormatBool(a.Compact)},
		"event":      {a.Event},
	}
	fullUrl := fmt.Sprintf(m.AnnounceURL, params.Encode())

	return fullUrl

}
