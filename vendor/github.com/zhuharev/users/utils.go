package users

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
)

func RandString(l int) string {
	buffer := make([]byte, l/2)
	if _, err := rand.Read(buffer); err != nil {
		panic(err)
	}

	return hex.EncodeToString(buffer)
}

func SecretHash(in, salt string) string {
	data := []byte(in)
	return fmt.Sprintf("%x", sha256.Sum256(data))
}

func ValidateSecretHash(hash, in, salt string) bool {
	return SecretHash(in, salt) == hash
}
