package client

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

const (
	chunkSize                    = 100 // In bytes
	discoveryServerBaseUrl       = "http://localhost:8081/api"
	discoveryServerWhereToUpload = discoveryServerBaseUrl + "/where-to-upload"
	discoveryServerIntroduceMe   = discoveryServerBaseUrl + "/introduce-me"
)

var logger *log.Logger

func encyptChunk(chunk []byte) []byte {
	// TODO: Do the encryption
	return chunk
}

func getChunkHash(chunk []byte) string {
	hashFunction := md5.New()
	hashFunction.Write(chunk)

	return hex.EncodeToString(hashFunction.Sum(nil))
}

func getUploadServers() []string {
	resp, err := http.Get(discoveryServerWhereToUpload)
	if err != nil {
		logger.Fatalln(err)
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)

	data := struct {
		Addresses []string
	}{}
	json.Unmarshal(body, &data)

	return data.Addresses
}

func askForIntroduction(serverId, hash string, chunkSize int) string {

	queryString := fmt.Sprintf("?to=%v&size=%v&hash=%v", serverId, chunkSize, hash)

	resp, err := http.Get(discoveryServerIntroduceMe + queryString)
	if err != nil {
		logger.Fatalln(err)
	}

	if resp.StatusCode != 200 {
		return ""
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)

	fmt.Println(body, resp.Status)
	data := struct {
		IntroductionKey string `json:"introduction-key"`
	}{}
	json.Unmarshal(body, &data)

	return data.IntroductionKey
}

func uploadChunk(chunk []byte) error {

	hash := getChunkHash(chunk)
	logger.Printf("Uploading chunk with md5 %v\n", hash)

	servers := getUploadServers()
	if len(servers) == 0 {
		logger.Fatalln("No avaliable servers to upload to.")
	}
	logger.Printf("This chunk should be uploaded to %v.\n", servers)

	for _, server := range servers {
		introductionKey := askForIntroduction(server, hash, len(chunk))
		logger.Printf("Got introduction key for server %v : %v.\n", server, introductionKey)
		if introductionKey == "" {
			logger.Fatalf("Failed to get the introduction key for the server %v.\n", server)
		}

		// Now it's time to start the actual upload... FARGHAAAAAAAAAAAAL!!

	}

	return nil
}

func Upload(filename string, l *log.Logger) {
	logger = l

	fileContent, err := ioutil.ReadFile(filename)
	if err != nil {
		logger.Fatalf("Failed to read the file : %v\n", err.Error())
	}

	logger.Printf("Successfully read the file (%v bytes).\n", len(fileContent))

	chunks := make([][]byte, 0)

	for i := 0; i < len(fileContent); i += chunkSize {
		if i+chunkSize > len(fileContent) {
			chunks = append(chunks, fileContent[i:])
		} else {
			chunks = append(chunks, fileContent[i:i+chunkSize])
		}
	}

	for _, chunk := range chunks {
		uploadChunk(encyptChunk(chunk))
	}

}
