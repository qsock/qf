package encrypt

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/pem"
)

func RsaSign(msg, Key []byte) (cryptText []byte, err error) {
	block, _ := pem.Decode(Key)

	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	myHash := sha256.New()
	myHash.Write(msg)
	hashed := myHash.Sum(nil)
	sign, err := rsa.SignPKCS1v15(rand.Reader, privateKey, crypto.SHA256, hashed)
	if err != nil {
		return nil, err
	}
	return sign, nil
}

func RsaVerifySign(msg []byte, sign []byte, Key []byte) bool {
	block, _ := pem.Decode(Key)

	publicInterface, _ := x509.ParsePKIXPublicKey(block.Bytes)
	publicKey := publicInterface.(*rsa.PublicKey)
	myHash := sha256.New()
	myHash.Write(msg)
	hashed := myHash.Sum(nil)
	result := rsa.VerifyPKCS1v15(publicKey, crypto.SHA256, hashed, sign)
	return result == nil
}
