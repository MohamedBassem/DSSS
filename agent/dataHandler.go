package agent

import (
	"io/ioutil"
	"os"
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


func HshHash(hash string) bool {
	_, err := os.Stat(datadir + hash)
	if err != nil {
		return true
	}
	return false
}
