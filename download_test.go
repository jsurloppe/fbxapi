package fbxapi

import (
	"testing"
)

func TestDownload(t *testing.T) {
	var data DownloadTask
	req := &DownloadReq{
		DownloadUrl: "http://ftp.free.fr/mirrors/cdimage.debian.org/debian-cd/current/i386/iso-cd/debian-9.3.0-i386-netinst.iso",
		DownloadDir: EncodePath("/Disque dur/Téléchargements/"),
	}

	err := testClient.Query(AddDownloadEP).WithFormBody(req).Do(&data)
	failOnError(t, err)
}
