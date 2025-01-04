package utils

import (
	"errors"
	"golang.org/x/crypto/bcrypt"
)

// HashData hashes the provided data using bcrypt.
// The saltRounds parameter defines the complexity (number of rounds).
// Only the hash is returned, as the salt is embedded within it.
func HashData(data string, saltRounds int) (string, error) {
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(data), saltRounds)
	if err != nil {
		return "", err
	}
	return string(hashedBytes), nil
}

// CompareHashes compares the "plain" data with previously hashed data.
// Returns true if they match.
func CompareHashes(pureData, hashedData string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(hashedData), []byte(pureData))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}
