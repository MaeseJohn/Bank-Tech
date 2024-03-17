package models

import (
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	UserId   string `validate:"required,uuid"`
	Email    string `validate:"required,email"`
	Password string `validate:"required,alphanum"`
	Name     string `validate:"required,alpha"`
	Funds    int    `validate:"required"`
	Role     string `validate:"required,oneof=issuer investor"`
}

func (u *User) HashSaltPassword() error {
	password := []byte(u.Password)
	//Hashing the password with the default cost of 10
	//The salt is automatically (and randomly) generated upon hashing a password
	hashedPassword, err := bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)
	u.Password = string(hashedPassword)
	return err
}

// Returns true if the password is correct, false if not
func (u *User) ValidatePassword(hashedPassword string, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}

// User constructor
func NewUser(userid, email, password, name, role string, funds int) *User {
	var user User
	user.UserId = userid
	user.Email = email
	user.Password = password
	user.Name = name
	user.Role = role
	user.Funds = funds

	return &user
}
