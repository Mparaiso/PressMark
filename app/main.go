package main
/*
import (
	"database/sql"
	_ "github.com/amattn/go-sqlite3"
	"log"
	"github.com/mparaiso/PressMark"
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
	user := &pressmark.User{Name: "John", Email: "John@yahoo.com"}
	userRepository := &pressmark.UserRepository{db}
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
	users := []*pressmark.User{
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
*/
