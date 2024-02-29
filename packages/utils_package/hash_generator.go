package utils_package

import (
	"crypto/sha1"
	"fmt"
)

func HashPassword(password string) string {
	h := sha1.New()
	h.Write([]byte(password))
	hashedPassword := fmt.Sprintf("%x", h.Sum(nil))
	return hashedPassword
}