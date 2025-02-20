package utility

import (
	"crypto/rand"
	"encoding/base32"
	"log"
)

func GetToken(length int) string {
	randomBytes := make([]byte, 32)
	_, err := rand.Read(randomBytes)
	if err != nil {
		log.Println("Error generating random bytes for token:", err)
	}
	return base32.StdEncoding.EncodeToString(randomBytes)[:length]
}
