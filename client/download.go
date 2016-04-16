package client

import (
	"bytes"
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

func decryptChunk(chunk string) (string, error) {
	return chunk, nil
}

func Download(manifestFileName, outputFileName string, l *log.Logger) {

	logger = l

	if manifestFileName == "" || outputFileName == "" {
		logger.Fatalln("Both the manifestFileName and the outputFileName should be specified")
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

		servers, err := whoHas(hash)
		if err != nil {
			logger.Fatalf("Failed to know who has %v with error %v\n", hash, err.Error())
		}
		encryptedChunk, err := downloadChunk(servers, hash)
		if err != nil {
			logger.Fatalf("Failed to download %v with error %v\n", hash, err.Error())
		}
		chunk, err := decryptChunk(encryptedChunk)
		if err != nil {
			logger.Fatalf("Failed to decrypt %v with error %v\n", hash, err.Error())
		}
		outputFile.WriteString(chunk)
		logger.Printf("%v successfully downloaded.\n", hash)
	}

	outputFile.Close()
	logger.Printf("The file was successfully downloaded to %v\n", outputFileName)
}
