package pressmark_test

import (
	"database/sql"
	_ "github.com/amattn/go-sqlite3"
	"github.com/interactiv/expect"
	"github.com/mparaiso/PressMark"
	"github.com/rubenv/sql-migrate"
	"testing"
)

type hash map[string]interface{}

func Before() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		return nil, err
	}
	migrations := &migrate.FileMigrationSource{
		Dir: "app/migrations/development.sqlite3",
	}
	_, err = migrate.Exec(db, "sqlite3", migrations, migrate.Up)
	if err != nil {
		return nil, err
	}
	return db, nil
}

func TestSave(t *testing.T) {
	e := expect.New(t)
	db, err := Before()
	e.Expect(err).ToBeNil()
	userRepository := &pressmark.UserRepository{DB: db}
	password := "the password"
	user := &pressmark.User{
		Name:           "John Doe",
		Email:          "john.doe@acme.com",
		SecurePassword: &pressmark.SecurePassword{Password: password},
	}
	err = userRepository.Save(user)
	e.Expect(err).ToBeNil()
	e.Expect(user.PasswordDigest).Not().ToBeNil()
	t.Log(user,user.PasswordDigest)
    
	new_name := "John Walker Doe"

	err = userRepository.UpdateAttribute(user, hash{"Name": new_name})
	e.Expect(err).ToBeNil()
	e.Expect(user.Name).ToBe(new_name)
   // t.Log(user.PasswordDigest)
}
