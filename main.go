package main

import (
	"database/sql"
	"fmt"
	_ "github.com/amattn/go-sqlite3"
	"log"
	"reflect"
)

import (
	"time"
)

// Article is a blog article
type Article struct {
	ID      int
	Title   string
	Content string
	Created time.Time
	Updated time.Time
}

// User is a user
type User struct {
	ID      int64
	Name    string
	Email   string
	Created time.Time
	Updated time.Time
}

func (user User) String() string {
	return fmt.Sprintf("{ID:%d,Name:%s}", user.ID, user.Name)
}

// UserRepository is a repository of users
type UserRepository struct {
	DB *sql.DB
}

// Find finds a user by id
func (repository *UserRepository) Find(id int64) (*User, error) {
	row := repository.DB.QueryRow("SELECT ID,NAME,EMAIL,CREATED,UPDATED FROM USERS WHERE ID=? LIMIT 1", id)
	user := &User{}
	err := row.Scan(&user.ID, &user.Name, &user.Email, &user.Created, &user.Updated)
	if err != nil {
		return nil, err
	}
	return user, err
}

// FindBy find users by fields
func (repository *UserRepository) FindBy(fields map[string]interface{}) ([]*User, error) {
	values := []interface{}{}
	whereExpression := ""
	for key, value := range fields {
		values = append(values, value)
		if whereExpression == "" {
			whereExpression = fmt.Sprintf("%s = ?", key)
		} else {
			whereExpression = fmt.Sprintf("%s AND %s = ? ", whereExpression, key)
		}
	}
	records, err := repository.DB.Query(fmt.Sprintf("SELECT ID,NAME,EMAIL,CREATED,UPDATED FROM USERS WHERE %s;", whereExpression), values...)
	if err != nil {
		return nil, err
	}
	users := []*User{}
	defer records.Close()
	for records.Next() {
		if err := records.Err(); err != nil {
			return nil, err
		}
		user := &User{}
		records.Scan(&user.ID, &user.Name, &user.Email, &user.Created, &user.Updated)
		users = append(users, user)
	}
	return users, nil
}

// DeleteAll deletes all models
func (repository *UserRepository) DeleteAll() error {
	result, err := repository.DB.Exec("DELETE FROM USERS;")
	if err != nil {
		return err
	}
	if _, err := result.RowsAffected(); err != nil {
		return err
	}
	return nil
}

// Save saves a model
func (repository *UserRepository) Save(user *User) error {
	if user.ID == 0 {
		user.Created = time.Now()
		user.Updated = time.Now()
		result, err := repository.DB.Exec("INSERT INTO USERS(NAME,EMAIL,CREATED,UPDATED) VALUES(?,?,?,?);", user.Name, user.Email, user.Created, user.Updated)
		if err != nil {
			return err
		}
		id, err := result.LastInsertId()
		if err != nil {
			return err
		}
		user.ID = id
		return nil
	}
	user.Updated = time.Now()

	result, err := repository.DB.Exec("UPDATE USERS SET NAME=?,EMAIL=?,UPDATED=? WHERE ID=?;", user.Name, user.Email, user.Updated, user.ID)
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return fmt.Errorf("User with ID %d does not exist.", user.ID)
	}
	return nil
}

// UpdateAttribute update selected attributes
func (repository *UserRepository) UpdateAttribute(user *User, attributes map[string]interface{}) error {
	values := []interface{}{}
	setStatement := ""
	userType := reflect.TypeOf(user)
	for key, value := range attributes {
		if _, ok := userType.Elem().FieldByName(key); !ok {
			return fmt.Errorf("type %s doesn't have a field named %s ", userType.Name(), key)
		}
		values = append(values, value)
		if setStatement == "" {
			setStatement = fmt.Sprintf(" %s = ?", key)
		} else {
			setStatement = fmt.Sprintf(" %s, %s = ?", setStatement, key)
		}
	}
	id := user.ID
	result, err := repository.DB.Exec(fmt.Sprintf("UPDATE USERS SET %s WHERE ID = ?;", setStatement), append(values, id)...)
	if err != nil {
		return err
	}
	if _, err := result.RowsAffected(); err != nil {
		return err
	}
	u, err := repository.Find(id)
	*user = *u
	//*user = *u
	if err != nil {
		return err
	}
	return nil
}

// Destroy removes a user
func (repository *UserRepository) Destroy(user *User) error {
	result, err := repository.DB.Exec("DELETE FROM USERS WHERE USERS.ID=?", user.ID)
	if err != nil {
		return err
	}
	if rows, err := result.RowsAffected(); err != nil {
		return err
	} else if rows <= 0 {
		return fmt.Errorf("User with ID %d could not be found and destroyed.", user.ID)
	}
	return nil
}

func main() {
	db, err := sql.Open("sqlite3", "pressmark.db")
	if err != nil {
		log.Fatal(err)
	}
	sqlRows, err := db.Query("SELECT ID,TITLE,CONTENT,CREATED,UPDATED FROM ARTICLES;")
	if err != nil {
		log.Fatal(err)
	}
	articles := []Article{}
	defer sqlRows.Close()
	for sqlRows.Next() == true {
		article := Article{}
		err := sqlRows.Scan(&article.ID, &article.Title, &article.Content,
			&article.Created, &article.Updated)
		if err != nil {
			log.Fatal(err)
		}
		articles = append(articles, article)
	}
	log.Print(articles)
	user := &User{Name: "John", Email: "John@yahoo.com"}
	userRepository := &UserRepository{db}
	userRepository.DeleteAll()
	err = userRepository.Save(user)
	if err != nil {
		log.Fatal(err)
	}
	user, err = userRepository.Find(user.ID)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("found user", user)
	err = userRepository.Destroy(user)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("destroyed user", user)
	users := []*User{
		{Name: "John", Email: "john@acme.com"},
		{Name: "Jane", Email: "jane@acme.com"},
	}
	for _, user := range users {
		err = userRepository.Save(user)
		if err != nil {
			log.Fatal(err)
		}
	}
	users, err = userRepository.FindBy(map[string]interface{}{"Name": "John"})
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("%v", users)
	err = userRepository.UpdateAttribute(users[0], map[string]interface{}{"Name": "Jack", "Email": "jack@acme.com"})
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("user updated by attribute: %s. \n", users[0])

}
