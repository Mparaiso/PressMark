package datamapper_test

import (
	"database/sql"
	"testing"

	_ "github.com/amattn/go-sqlite3"
	p "github.com/mparaiso/PressMark"
	"github.com/mparaiso/PressMark/datamapper"
	"github.com/rubenv/sql-migrate"
)

type hash map[string]interface{}

func Before(t *testing.T) *datamapper.DataMapper {
	db, err := sql.Open("sqlite3", ":memory:")

	if err != nil {
		t.Fatal(err)
	}
	migrations := &migrate.FileMigrationSource{
		Dir: "migrations/development.sqlite3",
	}
	_, err = migrate.Exec(db, "sqlite3", migrations, migrate.Up)
	if err != nil {
		t.Fatal(err)
	}

	connection := datamapper.NewConnectionWithOptions("sqlite3", db, &datamapper.ConnectionOptions{Logger: t})
	dm := datamapper.NewDataMapper(connection)
	err = dm.Register(&p.User{})
	if err != nil {
		t.Fatal(err)
	}
	return dm
}

func TestAll(t *testing.T) {
	dm := Before(t)
	userRepository, err := dm.GetRepository(&p.User{})
	if err != nil {
		t.Fatal(err)
	}
	user := &p.User{Name: "John Doe", Email: "john.doe@acme.com"}
	err = userRepository.Save(user)
	users := []*p.User{}
	err = userRepository.All(&users)
	// t.Log("users length : ", len(users))
	if err != nil {
		t.Fatal(err)
	}
	if len(users) != 1 {
		t.Fatalf("len(users) should be 1, got %d", len(users))
	}
}

func TestFind(t *testing.T) {
	dm := Before(t)
	userRepository, err := dm.GetRepository(&p.User{})
	if err != nil {
		t.Fatal(err)
	}
	user := &p.User{Name: "John Doe", Email: "john.doe@acme.com"}
	_ = user.GenerateSecurePassword("password")
	err = userRepository.Save(user)
	if err != nil {
		t.Fatal(err)
	}
	fetchedUser := &p.User{}
	err = userRepository.Find(user.ID, fetchedUser)
	if err != nil {
		t.Fatal(err)
	}
	// verify that the password is the right one
	err = fetchedUser.Authenticate("password")
	if err != nil {
		t.Fatal("Failed to authenticate.", err)
	}
}

func TestFindBy(t *testing.T) {
	dm := Before(t)
	userRepository, err := dm.GetRepository(&p.User{})
	if err != nil {
		t.Fatal(err)
	}
	userName := "John Doe"
	user := &p.User{Name: userName, Email: "john.doe@acme.com"}
	_ = user.GenerateSecurePassword("password")
	err = userRepository.Save(user)
	if err != nil {
		t.Fatal(err)
	}
	err = userRepository.Save(&p.User{Name: "Jane Doe", Email: "jane.doe@acme.com"})
	if err != nil {
		t.Fatal(err)
	}
	if id := user.ID; id == 0 {
		t.Fatal("user.ID should be >0, got", id)
	}
	candidates := []*p.User{}
	err = userRepository.FindBy(datamapper.Finder{Criteria: hash{"Name": userName}}, &candidates)
	if err != nil {
		t.Fatal(err)
	}
	if l := len(candidates); l != 1 {
		t.Fatal(candidates, "length should be 1, got", l)
	}
}
