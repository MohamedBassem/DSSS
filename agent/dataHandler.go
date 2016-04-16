package agent

import (
	"io/ioutil"
)

const (
	datadir	= "data/"
)


func Store(hash string, content string) {

	cnt := []byte(content)
	err := ioutil.WriteFile(datadir + hash, cnt, 0644)
	if err != nil {
		panic(err)
	}

}

