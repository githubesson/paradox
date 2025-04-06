package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"errors"
)

type ChromeDecryptionType int

const (
	ChromeCookie ChromeDecryptionType = iota
	ChromeLogin
)

func aes128CBCDecrypt(key, iv, encryptPass []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, errors.New("failed to create cipher")
	}
	encryptLen := len(encryptPass)

	dst := make([]byte, encryptLen)
	mode := cipher.NewCBCDecrypter(block, iv)
	mode.CryptBlocks(dst, encryptPass)
	dst = PKCS5UnPadding(dst)
	return dst, nil
}

func DecryptChromeAES(secretKey, encryptValue []byte, decryptType ChromeDecryptionType) ([]byte, error) {
	if len(secretKey) == 0 {
		return nil, errors.New("security key is empty")
	}

	if len(encryptValue) < 19 {
		return nil, errors.New("encrypted value too short")
	}

	chromeIV := []byte{32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32}

	value, err := aes128CBCDecrypt(secretKey, chromeIV, encryptValue[3:])
	if err != nil {
		return nil, err
	}

	return value, nil
}
