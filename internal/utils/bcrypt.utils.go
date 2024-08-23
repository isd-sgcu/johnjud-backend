package utils

import "golang.org/x/crypto/bcrypt"

type IBcryptUtil interface {
	GenerateHashedPassword(password string) (string, error)
	CompareHashedPassword(hashedPassword string, plainPassword string) error
}

type bcryptUtil struct{}

func NewBcryptUtil() IBcryptUtil {
	return &bcryptUtil{}
}

func (u *bcryptUtil) GenerateHashedPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hashedPassword), err
}

func (u *bcryptUtil) CompareHashedPassword(hashedPassword string, plainPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(plainPassword))
}
