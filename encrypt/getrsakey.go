package encrypt

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"os"
)

func GetRsaKey(privateKeyPrefix, privateFileName string,
	publicFileName, publicKeyPrefix string) error {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return err
	}
	x509PrivateKey := x509.MarshalPKCS1PrivateKey(privateKey)
	privateFile, err := os.Create(privateFileName)
	if err != nil {
		return err
	}
	defer privateFile.Close()
	privateBlock := pem.Block{
		Type:  privateKeyPrefix,
		Bytes: x509PrivateKey,
	}

	if err = pem.Encode(privateFile, &privateBlock); err != nil {
		return err
	}
	publicKey := privateKey.PublicKey
	x509PublicKey, err := x509.MarshalPKIXPublicKey(&publicKey)
	if err != nil {
		return err
	}
	publicFile, _ := os.Create(publicFileName)
	defer publicFile.Close()
	publicBlock := pem.Block{
		Type:  publicKeyPrefix,
		Bytes: x509PublicKey,
	}
	if err = pem.Encode(publicFile, &publicBlock); err != nil {
		return err
	}
	return nil
}
