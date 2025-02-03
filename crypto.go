package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
)

// Criptografa um arquivo
func encryptFile(filename string, key []byte) error {
	if filepath.Ext(filename) == ".enc" { // Evita criptografar arquivos já criptografados
		return nil
	}

	plaintext, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return err
	}

	ciphertext := gcm.Seal(nonce, nonce, plaintext, nil)

	err = ioutil.WriteFile(filename+".enc", ciphertext, 0644)
	if err != nil {
		return err
	}

	return os.Remove(filename) // Remove o arquivo original
}

// Descriptografa um arquivo
func decryptFile(filename string, key []byte) error {
	if filepath.Ext(filename) != ".enc" { // Evita tentar descriptografar arquivos que não são criptografados
		return nil
	}

	ciphertext, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return err
	}

	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		return fmt.Errorf("arquivo corrompido ou senha incorreta")
	}

	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]

	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return err
	}

	originalFilename := filename[:len(filename)-4]
	err = ioutil.WriteFile(originalFilename, plaintext, 0644)
	if err != nil {
		return err
	}

	return os.Remove(filename)
}

func processDirectory(directory string, key []byte, encrypt bool) error {
	return filepath.Walk(directory, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			if encrypt {
				return encryptFile(path, key)
			} else {
				if filepath.Ext(path) == ".enc" {
					return decryptFile(path, key)
				}
			}
		}
		return nil
	})
}
