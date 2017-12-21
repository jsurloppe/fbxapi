package fbxapi

type FTPConfig struct {
	Enabled             bool `json:"enabled"`
	AllowAnonymous      bool `json:"allow_anonymous"`
	AllowAnonymousWrite bool `json:"allow_anonymous_write"`
	WeakPassword        bool `json:"weak_password"`
	AllowRemoteAccess   bool `json:"allow_remote_access"`
	PortCtrl            int  `json:"port_ctrl"`
	PortData            int  `json:"port_data"`
}

var CurrentFTPConfigEP = &Endpoint{
	Verb: HTTP_METHOD_GET,
	Url:  "ftp/config/",
}

var UpdateFTPConfigEP = &Endpoint{
	Verb:         HTTP_METHOD_PUT,
	Url:          "ftp/config/",
	BodyRequired: true,
}
