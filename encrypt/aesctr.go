package encrypt

import (
	"crypto/aes"
	"crypto/cipher"
)

func AesCtrEncrypt(plainText, key, ivAes []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	stream := cipher.NewCTR(block, ivAes)
	cipherText := make([]byte, len(plainText))
	stream.XORKeyStream(cipherText, plainText)
	return cipherText, nil
}

func AesCtrDecrypt(cipherText, key, ivAes []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	stream := cipher.NewCTR(block, ivAes)
	plainText := make([]byte, len(cipherText))
	stream.XORKeyStream(plainText, cipherText)
	return plainText, nil
}
