package structs

type PingRequest struct{}

func (*PingRequest) String() string {
	return "PING"
}

type WhoHasRequest struct {
	Hash string
}

func (w *WhoHasRequest) String() string {
	return "WHO_HAS " + w.Hash
}

type UploadRequest struct {
	Hash    string
	Content string
}

func (u *UploadRequest) String() string {
	return "UPLOAD " + u.Hash + " " + u.Content
}

type DownloadRequest struct {
	Hash string
}

func (d *DownloadRequest) String() string {
	return "Download " + d.Hash
}
