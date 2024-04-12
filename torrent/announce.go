package torrent

func LoadTorrent(bval any) (MetaInfo, error) {
	m, _ := bval.(map[string]any)
	metainfo := MetaInfo{}
	metainfo.Announce = m["announce"].(string)
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
