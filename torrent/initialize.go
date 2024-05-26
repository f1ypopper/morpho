package torrent

import (
	"crypto/sha1"
	"io"
)

func CreateAnnounceData(metainfo *MetaInfo, info_dic map[string]any) AnnounceData {

	info := info_dic["info"].(map[string]any)
	raw_info := info["raw"].(string)
	h := sha1.New()
	io.WriteString(h, raw_info)
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
	if alist, ok := m["announce-list"]; ok {
		metainfo.AnnounceList = ManageAnnounceList(alist.([]interface{}))
	} else {
		metainfo.AnnounceList = ManageAnnounceList([]any{[]any{m["announce"].(string)}})
	}
	metainfo.Info, _ = loadInfo(m["info"])
	return metainfo, nil
}

func loadInfo(bval any) (TorrentInfo, error) {
	m, _ := bval.(map[string]any)
	tinfo := TorrentInfo{}
	tinfo.PieceLength = uint(m["piece length"].(int))
	hashes := []byte(m["pieces"].(string))
	num := len(m["pieces"].(string)) / 20
	for i := 0; i < num; i++ {
		var arr [20]byte
		copy(arr[:], hashes[i*20:(i+1)*20])
		tinfo.Pieces = append(tinfo.Pieces, arr)
	}
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
