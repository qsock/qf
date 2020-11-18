package encrypt

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
)

func RsaEncrypt(plainText, key []byte) (cryptText []byte, err error) {
	block, _ := pem.Decode(key)
	publicKeyInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	publicKey := publicKeyInterface.(*rsa.PublicKey)

	cipherText, err := rsa.EncryptPKCS1v15(rand.Reader, publicKey, plainText)
	if err != nil {
		return nil, err
	}
	return cipherText, nil
}

func RsaDecrypt(cryptText, key []byte) ([]byte, error) {
	block, _ := pem.Decode(key)

	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	plainText, err := rsa.DecryptPKCS1v15(rand.Reader, privateKey, cryptText)
	if err != nil {
		return nil, err
	}
	return plainText, nil
}
