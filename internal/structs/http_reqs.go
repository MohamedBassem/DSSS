package structs

type UploadRequestJSON struct {
	To      string
	Hash    string
	Content string
}

type DownloadRequestJSON struct {
	From string
	Hash string
}

type DownloadResponseJSON struct {
	Content string
	Hash    string
}

type WhoHasResponseJSON struct {
	Addresses []string
}
