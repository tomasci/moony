package crypto

import (
	"crypto/rand"
	"encoding/base64"
	"golang.org/x/crypto/argon2"
	"strings"
)

const saltSize = 16

// GenerateSalt – generates a random salt of a giver size
func GenerateSalt(size int) ([]byte, error) {
	// slice of bytes for salt of specified size
	salt := make([]byte, size)
	// fill slice with random data
	_, err := rand.Read(salt)

	if err != nil {
		return nil, err
	}

	return salt, nil
}

// HashCreate – generates hash of input with argon2
func HashCreate(input string) (string, error) {
	// generate salt
	salt, err := GenerateSalt(saltSize)
	if err != nil {
		return "", err
	}

	// create hash
	hash := argon2.IDKey([]byte(input), salt, 1, 64*1024, 4, 32) // default params from docs
	// concatenate hash and salt (salt will be used for verification)
	hashWithSalt := append(hash, salt...)

	// encode result as base64 string (for easy storage in database etc)
	return base64.RawStdEncoding.EncodeToString(hashWithSalt), nil
}

// HashValidate – check if input matches hash
func HashValidate(input string, encodedHash string) (bool, error) {
	// decode hash into byte slice
	hashWithSalt, err := base64.RawStdEncoding.DecodeString(encodedHash)
	if err != nil {
		return false, err
	}

	// get hash length by subtracting salt size
	hashLength := len(hashWithSalt) - saltSize
	// extract hash & salt
	hash := hashWithSalt[:hashLength]
	salt := hashWithSalt[hashLength:]

	// hash input same way as in HashCreate (but this time salt is not random, but from already existing hash)
	inputHash := argon2.IDKey([]byte(input), salt, 1, 64*1024, 4, 32)

	// compare and return result
	return strings.Compare(string(hash), string(inputHash)) == 0, nil
}
