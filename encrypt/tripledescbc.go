package encrypt

import (
	"crypto/cipher"
	"crypto/des"
)

func TripleDesEncrypt(plainText, key, ivDes []byte) ([]byte, error) {
	block, err := des.NewTripleDESCipher(key)
	if err != nil {
		return nil, err
	}
	paddingText := PKCS5Padding(plainText, block.BlockSize())
	blockMode := cipher.NewCBCEncrypter(block, ivDes)
	cipherText := make([]byte, len(paddingText))
	blockMode.CryptBlocks(cipherText, paddingText)
	return cipherText, nil
}

func TripleDesDecrypt(cipherText, key, ivDes []byte) ([]byte, error) {

	block, err := des.NewTripleDESCipher(key)
	if err != nil {
		return nil, err
	}

	blockMode := cipher.NewCBCDecrypter(block, ivDes)
	paddingText := make([]byte, len(cipherText))
	blockMode.CryptBlocks(paddingText, cipherText)

	plainText, err := PKCS5UnPadding(paddingText)
	if err != nil {
		return nil, err
	}
	return plainText, nil
}
