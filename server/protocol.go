package server

import "fmt"

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

type IntroductionRequest struct {
	Address string
	Size    int
	Hash    string
}

func (i *IntroductionRequest) String() string {
	return fmt.Sprintf("INTRODUCTION_REQUEST %v %v %v", i.Address, i.Size, i.Hash)
}
