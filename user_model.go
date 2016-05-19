package pressmark

import (
	"golang.org/x/crypto/bcrypt"
	"time"
  
)

// User is a user
type User struct {
	ID      int64
	Name    string
	Email   string
	Created time.Time
	Updated time.Time
	*SecurePassword
}


func (user *User) BeforeCreate() (err error) {
	user.Created = time.Now()
	return
}

func (user *User) BeforeSave() (err error) {
	if err = user.SecurePassword.BeforeSave(); err != nil {
		return err
	}
	user.Updated = time.Now()
	return
}

type SecurePassword struct {
	Password             string
	PasswordConfirmation string
	PasswordDigest       []byte
}

// BeforeSave does some work before saving
// see http://stackoverflow.com/questions/23259586/bcrypt-password-hashing-in-golang-compatible-with-node-js
// then   err = bcrypt.CompareHashAndPassword(hashedPassword, password)
func (sp *SecurePassword) BeforeSave()  error {
	if sp.Password != "" {
		passwordDigest, err := bcrypt.GenerateFromPassword([]byte(sp.Password), bcrypt.DefaultCost)
        if err!=nil{
            return err
        }
    sp.PasswordDigest = passwordDigest
        
	}
	return nil
}

// Authenticate return an error if the password and PasswordDigest do not match
func(sp *SecurePassword) Authenticate(password string)error{
   return bcrypt.CompareHashAndPassword(sp.PasswordDigest, []byte(password))
}
