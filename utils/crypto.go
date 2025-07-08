package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"fmt"
	"github.com/IntelXLabs-LLC/go-common/codecs"
	"io"
)

func Encrypt(key, text []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	ciphertext := make([]byte, aes.BlockSize+len(text))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return nil, err
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], text)

	return codecs.NewCodec().Base64().EncodeBytes(ciphertext), nil
}

func Decrypt(key, text []byte) ([]byte, error) {
	text, err := codecs.NewCodec().Base64().DecodeBytes(text)
	if err != nil {
		return nil, err
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	if len(text) < aes.BlockSize {
		return nil, fmt.Errorf("ciphertext is too short")
	}

	iv := text[:aes.BlockSize]
	text = text[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(text, text)

	return text, nil
}
