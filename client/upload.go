package client

import (
	"bytes"
	"crypto/md5"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/MohamedBassem/DSSS/internal/structs"
)

const (
	chunkSize                       = 100 // In bytes
	discoveryServerBaseURL          = "http://localhost:8081/api"
	discoveryServerWhereToUploadURL = discoveryServerBaseURL + "/where-to-upload"
	relayURL                        = discoveryServerBaseURL + "/relay"
	whoHasURL                       = discoveryServerBaseURL + "/who-has"
	downloadURL                     = discoveryServerBaseURL + "/download"
)

var logger *log.Logger

func encyptChunk(chunk []byte, pubKey *rsa.PublicKey) ([]byte, error) {
	tmp, err := rsa.EncryptOAEP(sha256.New(), rand.Reader, pubKey, chunk, []byte(EncryptionDecryptionLabel))
	if err != nil {
		return nil, err
	}
	str := base64.StdEncoding.EncodeToString(tmp)
	return []byte(str), nil
}

func getChunkHash(chunk []byte) string {
	hashFunction := md5.New()
	hashFunction.Write(chunk)

	return hex.EncodeToString(hashFunction.Sum(nil))
}

func getUploadServers() []string {
	resp, err := http.Get(discoveryServerWhereToUploadURL)
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

func uploadChunk(chunk []byte) (string, error) {

	hash := getChunkHash(chunk)
	logger.Printf("Uploading chunk with md5 %v\n", hash)

	servers := getUploadServers()
	if len(servers) == 0 {
		logger.Fatalln("No avaliable servers to upload to.")
	}
	logger.Printf("This chunk should be uploaded to %v.\n", servers)

	for _, server := range servers {
		var req = structs.UploadRequestJSON{
			To:      server,
			Hash:    hash,
			Content: string(chunk),
		}

		reqJson, _ := json.Marshal(&req)

		resp, err := (&http.Client{}).Post(relayURL, "application/json", bytes.NewReader(reqJson))

		if err != nil || resp.StatusCode != 200 {
			if err != nil {
				logger.Fatalln("Failed to upload to server : " + err.Error())
			} else {
				logger.Fatalln("Failed to upload to server : ", resp.StatusCode)
			}
		}

		if resp.Body != nil {
			resp.Body.Close()
		}
	}

	return hash, nil
}

func Upload(filename, outputManifestName, privateKeyFilePath string, l *log.Logger) {

	logger = l

	if filename == "" || outputManifestName == "" || privateKeyFilePath == "" {
		logger.Fatalln("The file to upload ,output manifest name and the privateKeyFilePath should be specified")
	}

	_, pubKey, err := ParsePrivateKey(privateKeyFilePath)
	if err != nil {
		logger.Fatalf("Failed to parse private key : %v\n", err.Error())
	}

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

	hashes := []string{}
	for _, chunk := range chunks {
		encChunk, err := encyptChunk(chunk, pubKey)
		fmt.Println(len(encChunk))

		if err != nil {
			logger.Fatalln(err)
		}
		hash, err := uploadChunk(encChunk)
		if err != nil {
			logger.Fatalln(err)
		}
		hashes = append(hashes, hash)
	}

	manifestFileContent := strings.Join(hashes, "\n")
	ioutil.WriteFile(outputManifestName, []byte(manifestFileContent), 0600)
	logger.Printf("Done uploading, manifest file dumped in %v.\n", outputManifestName)

}
