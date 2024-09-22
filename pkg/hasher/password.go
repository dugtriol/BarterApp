package hasher

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

//type PasswordHasher interface {
//	//Hash(password string) string
//	HashPassword(password string) (string, error)
//	CheckPassword(password string, hashedPassword string) error
//}
//
//type SHA1Hasher struct {
//	salt string
//}
//
//func NewSHA1Hasher(salt string) *SHA1Hasher {
//	return &SHA1Hasher{salt: salt}
//}
//
//func (h *SHA1Hasher) Hash(password string) string {
//	hash := sha1.New()
//	hash.Write([]byte(password))
//
//	return fmt.Sprintf("%x", hash.Sum([]byte(h.salt)))
//}

func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("failed to hash password: %w", err)
	}

	return string(hashedPassword), nil
}

func CheckPassword(password string, hashedPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}
