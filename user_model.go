package pressmark

import (
	"golang.org/x/crypto/bcrypt"
	"time"
)

// User is a user
type User struct {
	ID             int64
	Name           string
	Email          string
	Created        time.Time
	Updated        time.Time
	PasswordDigest string `db:"password_digest"`	
}

func (user *User) BeforeCreate() (err error) {
	user.Created = time.Now()
	return
}

func (user *User) BeforeSave() (err error) {
	user.Updated = time.Now()
	return
}

// BeforeSave does some work before saving
// see http://stackoverflow.com/questions/23259586/bcrypt-password-hashing-in-golang-compatible-with-node-js
func (user *User) GenerateSecurePassword(password string) error {
	passwordDigest, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.PasswordDigest = string(passwordDigest)
	return nil
}

// Authenticate return an error if the password and PasswordDigest do not match
func (user User) Authenticate(password string) error {
	return bcrypt.CompareHashAndPassword([]byte(user.PasswordDigest), []byte(password))
}
