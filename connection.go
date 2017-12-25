package fbxapi

type ConnectionStatus struct {
	State         string `json:"state"`
	Type          string `json:"type"`
	Media         string `json:"media"`
	Ipv4          string `json:"ipv4"`
	Ipv6          string `json:"ipv6"`
	RateUp        int    `json:"rate_up"`
	RateDown      int    `json:"rate_down"`
	BandwidthUp   int    `json:"bandwidth_up"`
	BandwidthDown int    `json:"bandwidth_down"`
	BytesUp       int    `json:"bytes_up"`
	BytesDown     int    `json:"bytes_down"`
	Ipv4PortRange [2]int `json:"ipv4_port_range"`
}

var ConnectionEP = &Endpoint{
	Verb: HTTP_METHOD_GET,
	Url:  "connection/",
}
