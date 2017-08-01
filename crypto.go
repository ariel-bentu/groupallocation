package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"io"
)

func AESencryptInit(base64passcode string) (key []byte, iv []byte) {
	password := make([]byte, 128)
	l, _ := base64.StdEncoding.Decode(password, []byte(base64passcode))
	key = password[:l]

	iv = make([]byte, aes.BlockSize)

	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		panic(err)
	}
	return key, iv
}

func AESencrypt(key []byte, iv []byte, content string) string {

	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}

	byteContent := []byte(content)

	ciphertext := make([]byte, aes.BlockSize+len(byteContent))
	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], byteContent)
	copy(ciphertext, iv)

	return base64.StdEncoding.EncodeToString(ciphertext)
}
