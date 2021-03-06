package agent

import (
	"io/ioutil"
	"os"
)

const (
	datadir = "data/"
)

func Store(hash string, content string) {

	cnt := []byte(content)
	err := ioutil.WriteFile(datadir+hash, cnt, 0644)
	if err != nil {
		panic(err)
	}

}

func HasHash(hash string) bool {
	_, err := os.Stat(datadir + hash)
	if err != nil {
		return false
	}
	return true
}

func Fetch(hash string) string {
	exists := HasHash(hash)
	if !exists {
		return "ERROR Hash not found"
	}

	cnt, err := ioutil.ReadFile(datadir + hash)
	if err != nil {
		panic(err)
	}
	return string(cnt)

}
