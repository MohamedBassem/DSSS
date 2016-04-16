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

func Fetch(hash string) string {
	exists := HshHash(hash)
	if !exists {
		return "ERROR Hash not found"
	}
	
	cnt, err := ioutil.ReadFile(datadir + hash)
	if err != nil { panic(err) }
	return string(cnt)

}
