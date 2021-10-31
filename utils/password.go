package utils

import (
	"golang.org/x/crypto/bcrypt"
)

func HashAndSalt(pwd string) ([]byte, error) {
	// Use GenerateFromPassword to hash & salt pwd.
	// MinCost is just an integer constant provided by the bcrypt
	// package along with DefaultCost & MaxCost.
	// The cost can be any value you want provided it isn't lower
	// than the MinCost (4)
	hash, err := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	return hash, nil
}

func ComparePasswords(hashedPwd []byte, plainPwd []byte) bool {
	err := bcrypt.CompareHashAndPassword(hashedPwd, plainPwd)
	return err == nil
}
