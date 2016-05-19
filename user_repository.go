package pressmark

import (
	"database/sql"
	"fmt"
	"reflect"
)


// UserRepository is a repository of users
type UserRepository struct {
	DB *sql.DB
}

// Find finds a user by id
func (repository *UserRepository) Find(id int64) (*User, error) {
	row := repository.DB.QueryRow("SELECT ID,NAME,EMAIL,CREATED,UPDATED,PASSWORD_DIGEST FROM USERS WHERE ID=? LIMIT 1", id)
	user := &User{SecurePassword:&SecurePassword{}}
	err := row.Scan(&user.ID, &user.Name, &user.Email, &user.Created, &user.Updated,&user.PasswordDigest)
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
	if u, ok := interface{}(user).(BeforeSaveCallback); ok == true {
		if err := u.BeforeSave(); err != nil {
			return err
		}
	}
	if user.ID == 0 {
		if u, ok := interface{}(user).(BeforeCreateCallback); ok == true {
			if err := u.BeforeCreate(); err != nil {
				return err
			}
		}
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
	if u, ok := interface{}(user).(BeforeUpdateCallback); ok == true {
		if err := u.BeforeUpdate(); err != nil {
			return err
		}
	}
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
