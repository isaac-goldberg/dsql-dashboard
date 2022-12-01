package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/rand"
	"encoding/hex"
	"io"
	"os"
	"strconv"
	"strings"
)

func createHash(key string) string {
	hasher := md5.New()
	hasher.Write([]byte(key))
	return hex.EncodeToString(hasher.Sum(nil))
}

func Encrypt(d string) string {
	var passphrase string = os.Getenv("ENCRYPTION_PASSPHRASE")

	data := []byte(d)
	block, _ := aes.NewCipher([]byte(createHash(passphrase)))
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		panic(err.Error())
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		panic(err.Error())
	}
	ciphertext := gcm.Seal(nonce, nonce, data, nil)

	var str = ""
	for i, b := range ciphertext {
		str += strconv.Itoa(int(b))
		if i < len(ciphertext)-1 {
			str += ","
		}
	}
	return str
}

func Decrypt(d string) (string, error) {
	var passphrase string = os.Getenv("ENCRYPTION_PASSPHRASE")

	data := make([]byte, 0)
	split := strings.Split(d, ",")
	for _, b := range split {
		num, err := strconv.Atoi(b)
		if err != nil {
			return "", err
		}
		data = append(data, byte(num))
	}

	key := []byte(createHash(passphrase))
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}
	nonceSize := gcm.NonceSize()
	nonce, ciphertext := data[:nonceSize], data[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", err
	}
	return string(plaintext), nil
}
