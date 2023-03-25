package utils

import "golang.org/x/crypto/bcrypt"

func BcryptMake(str []byte) string {
	bytes, _ := bcrypt.GenerateFromPassword(str, 14)
	return string(bytes)
}

func BcryptMakeCheck(str, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(str))
	return err == nil
}
