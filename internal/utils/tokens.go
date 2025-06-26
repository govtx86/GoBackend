package utils

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
)

func GenerateToken(length int) (string, error) {
	bytes := make([]byte, length)
	_, err := rand.Read(bytes); 
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(bytes), nil
}

func HashToken(plainText string) string {
	hash := sha256.Sum256([]byte(plainText))
	return base64.URLEncoding.EncodeToString(hash[:])
}

func CheckTokenHash(plainText string, hash string) bool {
	return hash == HashToken(plainText)
}