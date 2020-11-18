package encrypt

import (
	"crypto/cipher"
	"crypto/des"
)

func DesCbcEncrypt(plainText, key, ivDes []byte) ([]byte, error) {
	block, err := des.NewCipher(key)
	if err != nil {
		return nil, err
	}
	paddingText := PKCS5Padding(plainText, block.BlockSize())
	blockMode := cipher.NewCBCEncrypter(block, ivDes)
	cipherText := make([]byte, len(paddingText))
	blockMode.CryptBlocks(cipherText, paddingText)
	return cipherText, nil
}

func DesCbcDecrypt(cipherText, key, ivDes []byte) ([]byte, error) {
	block, err := des.NewCipher(key)
	if err != nil {
		return nil, err
	}

	blockMode := cipher.NewCBCDecrypter(block, ivDes)
	plainText := make([]byte, len(cipherText))
	blockMode.CryptBlocks(plainText, cipherText)

	unPaddingText, err := PKCS5UnPadding(plainText)
	if err != nil {
		return nil, err
	}
	return unPaddingText, nil
}
