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
