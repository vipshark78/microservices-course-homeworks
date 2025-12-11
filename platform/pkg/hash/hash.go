package hash

import "golang.org/x/crypto/bcrypt"

type Hasher interface {
	Hash(text string) (string, error)
	Compare(hash, text string) error
}

func Hash(text string) (string, error) {
	b, err := bcrypt.GenerateFromPassword([]byte(text), bcrypt.MinCost)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func Compare(hash, text string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(text))
	return err
}
