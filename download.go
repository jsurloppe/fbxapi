package fbxapi

type Download struct {
	ID              int    `json:"id"`
	Type            string `json:"type"`
	Name            string `json:"name"`
	Status          string `json:"status"`
	Size            int    `json:"size"`
	QueuePos        int    `json:"queue_pos"`
	IOPriority      string `json:"io_priority"`
	TXBytes         int    `json:"tx_bytes"`
	RXBytes         int    `json:"rx_bytes"`
	TXRate          int    `json:"tx_rate"`
	RXRate          int    `json:"rx_rate"`
	TXPct           int    `json:"tx_pct"`
	RXPct           int    `json:"rx_pct"`
	Error           string `json:"error"`
	CreatedTS       int    `json:"created_ts"`
	ETA             int    `json:"eta"`
	DownloadDir     string `json:"download_dir"`
	StopRatio       int    `json:"stop_ratio"`
	ArchivePassword string `json:"stop_ratio"`
	InfoHash        string `json:"info_hash"`
	PieceLength     int    `json:"piece_length"`
}

type DownloadReq struct {
	DownloadUrl     string `json:"download_url,omitempty"`
	DownloadUrlList string `json:"download_url_list,omitempty"`
	DownloadDir     string `json:"download_dir"`
	Recursive       bool   `json:"recursive"`
	Username        string `json:"username,omitempty"`
	Password        string `json:"password,omitempty"`
	ArchivePassword string `json:"archive_password,omitempty"`
	Cookies         string `json:"cookies,omitempty"`
}

type DownloadTask struct {
	ID int `json:"id"`
}

// DownloadsEP endpoint definition
// Output: []Download
var DownloadsEP = &Endpoint{
	Verb: HTTP_METHOD_GET,
	Url:  "downloads/",
}

// DeleteDownloadEP endpoint definition
var DeleteDownloadEP = &Endpoint{
	Verb: HTTP_METHOD_DELETE,
	Url:  "downloads/{{.id}}/",
}

// EraseDownloadEP endpoint definition
var EraseDownloadEP = &Endpoint{
	Verb: HTTP_METHOD_DELETE,
	Url:  "downloads/{{.id}}/erase",
}

// AddDownloadEP endpoint definition
// Output: Download
var AddDownloadEP = &Endpoint{
	Verb: HTTP_METHOD_POST,
	Url:  "downloads/add/",
}
