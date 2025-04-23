package crypto

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/des"
	"crypto/sha1"
	"crypto/sha256"
	"database/sql"
	"encoding/asn1"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"paradox_server/models"
	"paradox_server/static"

	"golang.org/x/crypto/pbkdf2"
)

func GetHashedSalt(globalSalt []byte) []byte {
	hash := sha1.New()
	hash.Write(globalSalt)
	hash.Write([]byte{})
	return hash.Sum(nil)
}

func generatePBKDF2Key(hashedSalt []byte, params models.PBKDF2Params) []byte {
	return pbkdf2.Key(
		hashedSalt,
		params.Salt,
		int(params.IterationCount),
		int(params.KeyLength),
		sha256.New,
	)
}

func decryptAES(encryptedText, key, iv []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	blockMode := cipher.NewCBCDecrypter(block, iv)
	decryptedText := make([]byte, len(encryptedText))
	blockMode.CryptBlocks(decryptedText, encryptedText)
	return PKCS5UnPadding(decryptedText), nil
}

func decryptPBE(decodedItem models.EncryptedData, globalSalt []byte) ([]byte, error) {
	if Equal(decodedItem.EncryptionAlgo.AlgoID, static.Pkcs5Id) {
		params := decodedItem.EncryptionAlgo.Params.KDF.Params
		hashedSalt := GetHashedSalt(globalSalt)

		key := generatePBKDF2Key(hashedSalt, params)
		iv := append([]byte{0x04, 0x0e}, decodedItem.EncryptionAlgo.Params.Cipher.IV...)

		decryptedText, err := decryptAES(decodedItem.Encrypted, key, iv)
		if err != nil {
			return nil, err
		}

		return decryptedText, nil
	}

	return nil, errors.New("decryptPBE: unsupported algorithm")
}

func GetDecryptionKey(dbPath string) ([]byte, error) {
	fmt.Println("Db path: ", dbPath)

	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, errors.New("getDecryptionKey: cannot open db")
	}
	defer db.Close()

	var globalSalt []byte
	err = db.QueryRow("SELECT item1 FROM metadata WHERE id = 'password';").Scan(&globalSalt)
	if err != nil {
		return nil, errors.New("getDecryptionKey: cannot get global salt")
	}

	var a11, ckaIdValue []byte
	if err = db.QueryRow("SELECT a11,a102 FROM nssPrivate WHERE a11 IS NOT NULL LIMIT 1;").Scan(&a11, &ckaIdValue); err != nil {
		return nil, err
	}

	if !bytes.Equal(ckaIdValue, static.CkaId) {
		return nil, errors.New("getDecryptionKey: unsupported algorythm")
	}

	var decodedA11 models.EncryptedData
	if _, err := asn1.Unmarshal(a11, &decodedA11); err != nil {
		return nil, err
	}

	clearText, err := decryptPBE(decodedA11, globalSalt)
	if err != nil {
		return nil, err
	}
	return clearText[:24], nil
}

func decryptData(keyId, iv, ciphertext []byte, key []byte) ([]byte, error) {
	if !bytes.Equal(keyId, static.CkaId) {
		return nil, errors.New("decryptData: key ID does not match ckaId")
	}

	block, err := des.NewTripleDESCipher(key)
	if err != nil {
		return nil, errors.New("decryptData: error creating DES block")
	}

	blockMode := cipher.NewCBCDecrypter(block, iv)
	decrypted := make([]byte, len(ciphertext))
	blockMode.CryptBlocks(decrypted, ciphertext)
	return decrypted, nil

}

func DecryptCredentials(key []byte, login *models.Login) (*models.Login, error) {
	keyId, iv, ciphertext, err := decodeLoginData(login.EncryptedUsername)
	if err != nil {
		return nil, errors.New("decryptCredentials: error decoding username")
	}

	decryptedUsername, err := decryptData(keyId, iv, ciphertext, key)
	if err != nil {
		return nil, errors.New("decryptCredentials: error decrypting username")
	}

	keyId, iv, ciphertext, err = decodeLoginData(login.EncryptedPassword)
	if err != nil {
		return nil, errors.New("decryptCredentials: error decoding password")
	}

	decryptedPassword, err := decryptData(keyId, iv, ciphertext, key)
	if err != nil {
		return nil, errors.New("decryptCredentials: error decrypting password")
	}

	login.Username = string(PKCS5UnPadding(decryptedUsername))
	login.Password = string(PKCS5UnPadding(decryptedPassword))

	return login, nil
}

func decodeLoginData(data string) (keyID, iv, ciphertext []byte, err error) {
	decodedData, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		return nil, nil, nil, err
	}

	var asn1Data struct {
		KeyID      []byte              `asn1:""`
		Cypher     models.CipherParams `asn1:""`
		Ciphertext []byte              `asn1:""`
	}

	_, err = asn1.Unmarshal(decodedData, &asn1Data)
	if err != nil {
		return nil, nil, nil, errors.New("decodeLoginData: failed to unmarshal ASN.1 data")
	}

	return asn1Data.KeyID, asn1Data.Cypher.IV, asn1Data.Ciphertext, nil
}

func GetLoginsData(filepath string) ([]*models.Login, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return nil, errors.New("getLoginsData: cannot open the file")
	}
	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		return nil, errors.New("getLoginsData: file does not exist")
	}

	fileContent := make([]byte, fileInfo.Size())

	if _, err = file.Read(fileContent); err != nil {
		return nil, errors.New("getLoginsData: error reading file")
	}

	var data models.GeckoJsonData
	if err = json.Unmarshal(fileContent, &data); err != nil {
		return nil, errors.New("getLoginsData: error unmarshaling JSON")
	}

	var logins []*models.Login
	for _, jLogin := range data.Logins {
		login := &models.Login{
			EncryptedUsername: jLogin.EncryptedUsername,
			EncryptedPassword: jLogin.EncryptedPassword,
			URL:               jLogin.Hostname,
		}
		logins = append(logins, login)
	}

	return logins, nil
}
