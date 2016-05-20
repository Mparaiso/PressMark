package pressmark_test

import (
	"database/sql"
	_ "github.com/amattn/go-sqlite3"
	"github.com/jmoiron/sqlx"
	p "github.com/mparaiso/PressMark"
	"github.com/rubenv/sql-migrate"
	"testing"
)

type hash map[string]interface{}

func Before() (*sqlx.DB, error) {
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

	return sqlx.NewDb(db, "sqlite3"), nil
}

// func TestSave(t *testing.T) {
// 	e := expect.New(t)
// 	db, err := Before()
// 	e.Expect(err).ToBeNil()
// 	userRepository := &p.UserRepository{DB: db}
// 	user := &p.User{Name: "John Doe", Email: "john.doe@acme.com", SecurePassword: &p.SecurePassword{}}
// 	err = user.GenerateSecurePassword("password")
// 	if err != nil {
// 		t.Error(err)
// 	}
// 	err = userRepository.Save(user)
// 	if err != nil {
// 		t.Error(err)
// 	}
// 	t.Log(user)
// 	db.Close()
// }

func TestAll(t *testing.T) {
	db, _ := Before()
	defer db.Close()
	userRepository := &p.UserRepository{DB: db, TableName: "USERS",IDField:"ID"}
	user := &p.User{Name: "John Doe", Email: "john.doe@acme.com"}
	err := userRepository.Save(user)
	users := []*p.User{}
	err = userRepository.All(&users)
	t.Log("users length : ",len(users))
	if err != nil {
		t.Fatal(err)
	}
}

func TestFind(t *testing.T) {
	db, _ := Before()
	defer db.Close()
	userRepository := &p.UserRepository{DB: db, IDField: "ID", TableName: "USERS"}
	user := &p.User{Name: "John Doe", Email: "john.doe@acme.com"}
	_ = user.GenerateSecurePassword("password")
	err := userRepository.Save(user)
	var id int64 = 1 // user.ID
	t.Log("ID",id)
	fetchedUser := &p.User{}
	err = userRepository.Find(id, fetchedUser)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(fetchedUser)
	t.Log(fetchedUser.PasswordDigest)
	// verify that the password is the right one
	err = fetchedUser.Authenticate("password")
	if err != nil {
		t.Fatal("Failed to authenticate.", err)
	}

}
