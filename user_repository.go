package pressmark

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"reflect"
	"strings"
)

// UserRepository is a repository of users
type UserRepository struct {
	DB        *sqlx.DB
	IDField   string
	TableName string
}

// All finds all
func (repository *UserRepository) All(collection interface{}) error {
	return repository.DB.Select(collection, fmt.Sprintf("SELECT %[1]s.* FROM %[1]s ", repository.TableName))

}

// Find finds a user by id
func (repository *UserRepository) Find(id int64, model interface{}) error {
	return repository.DB.Get(model, fmt.Sprintf("SELECT %[1]s.* FROM %[1]s WHERE %[2]s =? ", repository.TableName, repository.IDField), id)
}

// FindBy find users by fields
func (repository *UserRepository) FindBy(fields map[string]interface{}, users []*User) error {
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
	err := repository.DB.Select(&users, fmt.Sprintf("SELECT * WHERE %s;", whereExpression), values...)
	if err != nil {
		return err
	}
	return nil
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
func (repository *UserRepository) Save(model *User) error {
	if u, ok := interface{}(model).(BeforeSaveCallback); ok == true {
		if err := u.BeforeSave(); err != nil {
			return err
		}
	}
	if model.ID == 0 {
		if u, ok := interface{}(model).(BeforeCreateCallback); ok == true {
			if err := u.BeforeCreate(); err != nil {
				return err
			}
		}
		paths := []string{}
		values := []interface{}{}
		fieldMap := repository.DB.Mapper.FieldMap(reflect.ValueOf(model))

		for key, value := range fieldMap {
			if strings.ToLower(key) != strings.ToLower(repository.IDField) {
				paths = append(paths, key)
				values = append(values, value)
			}
		}
		//fmt.Print(paths, values)
		query := fmt.Sprintf("INSERT INTO USERS(%s) VALUES(%s);",
			strings.Join(paths, ","),
			strings.Join(
				strings.Split(strings.Repeat("?", len(paths)), ""), ","))
		fmt.Print(query)
		result, err := repository.DB.Exec(query, values...)
		if err != nil {
			return err
		}
		id, err := result.LastInsertId()
		if err != nil {
			return err
		}
		model.ID = id
		return nil
	}
	if u, ok := interface{}(model).(BeforeUpdateCallback); ok == true {
		if err := u.BeforeUpdate(); err != nil {
			return err
		}
	}
	result, err := repository.DB.Exec("UPDATE USERS SET NAME=?,EMAIL=?,UPDATED=? WHERE ID=?;",
		model.Name, model.Email, model.Updated, model.ID)
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return fmt.Errorf("User with ID %d does not exist.", model.ID)
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
	err = repository.Find(id, user)
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
