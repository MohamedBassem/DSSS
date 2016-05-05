package client

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/MohamedBassem/DSSS/internal/structs"
)

func whoHas(hash string) ([]string, error) {
	resp, err := (&http.Client{}).Get(whoHasURL + "?q=" + hash)
	if err != nil {
		return nil, err
	}
	if resp.Body == nil {
		return nil, errors.New("Empty Body")
	}
	defer resp.Body.Close()

	var servers structs.WhoHasResponseJSON
	err = json.NewDecoder(resp.Body).Decode(&servers)
	if err != nil {
		return nil, err
	}
	return servers.Addresses, nil
}

func downloadChunk(servers []string, hash string) (string, error) {

	// TODO : TRY ALL SERVERS
	var req = structs.DownloadRequestJSON{
		From: servers[0],
		Hash: hash,
	}

	reqJson, _ := json.Marshal(&req)

	resp, err := (&http.Client{}).Post(downloadURL, "application/json", bytes.NewReader(reqJson))

	if err != nil {
		return "", err
	}

	if resp.StatusCode != 200 {
		return "", errors.New("Non 200 response : " + strconv.Itoa(resp.StatusCode))
	}

	defer resp.Body.Close()

	var respBody structs.DownloadResponseJSON
	err = json.NewDecoder(resp.Body).Decode(&respBody)
	if err != nil {
		return "", err
	}

	return respBody.Content, nil
}

func decryptChunk(chunk string, privKey *rsa.PrivateKey) (string, error) {
	tmp, err := base64.StdEncoding.DecodeString(chunk)
	if err != nil {
		return "", err
	}
	plain, err := rsa.DecryptOAEP(sha256.New(), rand.Reader, privKey, tmp, []byte(EncryptionDecryptionLabel))
	if err != nil {
		return "", err
	}
	return string(plain), nil
}

func Download(manifestFileName, outputFileName, privateKeyFilePath string, l *log.Logger) {

	logger = l

	if manifestFileName == "" || outputFileName == "" || privateKeyFilePath == "" {
		logger.Fatalln("The manifestFileName, outputFileName and privateKeyFilePath should be specified")
	}

	privKey, _, err := ParsePrivateKey(privateKeyFilePath)
	if err != nil {
		logger.Fatalf("Failed to parse private key : %v\n", err.Error())
	}

	manifestFileContent, err := ioutil.ReadFile(manifestFileName)
	if err != nil {
		logger.Fatalf("Failed to read the file : %v\n", err.Error())
	}

	chunkHashes := strings.Split(string(manifestFileContent), "\n")

	outputFile, err := os.Create(outputFileName)
	if err != nil {
		logger.Fatalf("Failed to open %v for writing with error %v\n", outputFileName, err.Error())
	}

	for _, hash := range chunkHashes {
		logger.Printf("Trying to download chunk with hash %v.\n", hash)
		servers, err := whoHas(hash)
		logger.Printf("Chunk with hash %v is on %v.\n", hash, servers)
		if err != nil {
			logger.Fatalf("Failed to know who has %v with error %v\n", hash, err.Error())
		}
		if len(servers) == 0 {
			logger.Fatalf("Failed to download %v no online servers having this chunk\n", hash)
		}
		encryptedChunk, err := downloadChunk(servers, hash)
		logger.Printf("Chunk with hash %v downloaded.\n", hash)
		if err != nil {
			logger.Fatalf("Failed to download %v with error %v\n", hash, err.Error())
		}
		chunk, err := decryptChunk(encryptedChunk, privKey)
		if err != nil {
			logger.Fatalf("Failed to decrypt %v with error %v\n", hash, err.Error())
		}
		outputFile.WriteString(chunk)
		logger.Printf("%v successfully downloaded.\n", hash)
	}

	outputFile.Close()
	logger.Printf("The file was successfully downloaded to %v\n", outputFileName)
}
