package utils

import (
	"crypto/sha256"
	"encoding/hex"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

var cost = 12

const (
	bcryptPrefix = "bcrypt$"
	sha256Prefix = "sha256$"
)

func Hash(value string) (string, error) {

	if len(value) > 72 {
		sum := sha256.Sum256([]byte(value))
		return sha256Prefix + hex.EncodeToString(sum[:]), nil
	}

	hash, err := bcrypt.GenerateFromPassword(
		[]byte(value),
		cost,
	)
	if err != nil {
		return "", err
	}

	return bcryptPrefix + string(hash), nil
}

func Compare(hash string, value string) error {
	switch {
	case strings.HasPrefix(hash, bcryptPrefix):
		return bcrypt.CompareHashAndPassword(
			[]byte(strings.TrimPrefix(hash, bcryptPrefix)),
			[]byte(value),
		)

	case strings.HasPrefix(hash, sha256Prefix):
		sum := sha256.Sum256([]byte(value))
		currentHash := hex.EncodeToString(sum[:])

		if currentHash != strings.TrimPrefix(hash, sha256Prefix) {
			return bcrypt.ErrMismatchedHashAndPassword
		}

		return nil

	default:
		return bcrypt.CompareHashAndPassword(
			[]byte(hash),
			[]byte(value),
		)
	}
}
