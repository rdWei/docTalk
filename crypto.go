package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha256"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
)

func deriveKey(password string) []byte {
	hash := sha256.New()
	hash.Write([]byte(password))
	return hash.Sum(nil)
}

func encrypt(plaintext, password string) (string, error) {
	key := deriveKey(password)

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	plaintextBytes := []byte(plaintext)
	blockSize := block.BlockSize()
	padding := blockSize - len(plaintextBytes)%blockSize
	paddingText := make([]byte, padding)
	plaintextBytes = append(plaintextBytes, paddingText...)

	ciphertext := make([]byte, len(plaintextBytes))
	iv := make([]byte, blockSize) 
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", err
	}

	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(ciphertext, plaintextBytes)

	ciphertext = append(iv, ciphertext...)
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

func decrypt(ciphertextBase64, password string) (string, error) {
	key := deriveKey(password)

	ciphertext, err := base64.StdEncoding.DecodeString(ciphertextBase64)
	if err != nil {
		return "", err
	}

	blockSize := aes.BlockSize
	if len(ciphertext) < blockSize {
		return "", fmt.Errorf("ciphertext too short")
	}

	iv := ciphertext[:blockSize]
	ciphertext = ciphertext[blockSize:]

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	mode := cipher.NewCBCDecrypter(block, iv)
	plaintextBytes := make([]byte, len(ciphertext))
	mode.CryptBlocks(plaintextBytes, ciphertext)

	padding := int(plaintextBytes[len(plaintextBytes)-1])
	plaintextBytes = plaintextBytes[:len(plaintextBytes)-padding]

	return string(plaintextBytes), nil
}
