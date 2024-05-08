package helpers

import "golang.org/x/crypto/bcrypt"

func Hash(plainText string) (string, error) {
	hashedText, errHash := bcrypt.GenerateFromPassword([]byte(plainText), bcrypt.DefaultCost)
	if errHash != nil {
		return "", errHash
	}

	return string(hashedText), nil
}

func VerifyHash(hashedText string, plainText string) error {
	errHash := bcrypt.CompareHashAndPassword([]byte(hashedText), []byte(plainText))
	if errHash != nil {
		return errHash
	}

	return nil
}
