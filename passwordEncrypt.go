package main

import (
	"golang.org/x/crypto/bcrypt"
)

func encryptPass(password string) ([]byte, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	return hashedPassword, nil
}

func verify(hashpass []byte, pass string) bool {
	err := bcrypt.CompareHashAndPassword(hashpass, []byte(pass))
	return err == nil
}