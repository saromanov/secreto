package crypto

// Попробовать сделать монолит а внутри микросервисы
// внутри несколько серверов

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"

	"github.com/pkg/errors"
)

// Encrypt provides encryption of the key
func Encrypt(stringToEncrypt string, secret string) (string, error) {
	if secret == "" {
		return "", errors.New("secret key is not defined")
	}
	if stringToEncrypt == "" {
		return "", errors.New("encrypt string is not defined")
	}
	
	plaintext := []byte(stringToEncrypt)

	block, err := aes.NewCipher([]byte(secret))
	if err != nil {
		return "", errors.Wrap(err, "unable to apply new cipher")
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", errors.Wrap(err, "unable to apply new gcm")
	}

	nonce := make([]byte, aesGCM.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return "", errors.Wrap(err, "unable to read data")
	}

	ciphertext := aesGCM.Seal(nonce, nonce, plaintext, nil)
	return fmt.Sprintf("%x", ciphertext), nil
}

// Decrypt provides decryption of the key
func Decrypt(encryptedString string, secret string) (string, error) {
	enc, err := hex.DecodeString(encryptedString)

	block, err := aes.NewCipher([]byte(secret))
	if err != nil {
		return "", errors.Wrap(err, "unable to apply new cipher")
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", errors.Wrap(err, "unable to apply new gcm")
	}

	nonceSize := aesGCM.NonceSize()
	nonce, ciphertext := enc[:nonceSize], enc[nonceSize:]

	plaintext, err := aesGCM.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", errors.Wrap(err, "unable to decode data")
	}

	return fmt.Sprintf("%s", plaintext), nil
}
