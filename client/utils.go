package client

import (
	"bufio"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"os"
)

const EncryptionDecryptionLabel = "TEST"

func ParsePrivateKey(privateKeyFilePath string) (*rsa.PrivateKey, *rsa.PublicKey, error) {
	file, _ := ioutil.ReadFile(privateKeyFilePath)

	pemdata, _ := pem.Decode(file)
	x, err := x509.DecryptPEMBlock(pemdata, []byte(""))
	if err != nil {
		fmt.Print("Please enter your private key's password: ")
		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()
		x, err = x509.DecryptPEMBlock(pemdata, []byte(scanner.Text()))
		if err != nil {
			return nil, nil, err
		}
	}

	pk, _ := x509.ParsePKCS1PrivateKey(x)
	pubkey := pk.Public().(*rsa.PublicKey)

	return pk, pubkey, nil
}
