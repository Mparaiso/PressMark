package datamapper

import (
	"fmt"
	"log"
	"reflect"
	"strings"
)

type Finder struct {
	Criteria map[string]interface{}
	OrderBy  map[string]string
	Limit    int64
	Offset   int64
}

func (finder Finder) AcceptRepository(repository *Repository) (string, []interface{}) {
	values := []interface{}{}
	whereStatement := ""
	orderByStatement := ""
	limitStatement := ""
	offsetStatement := ""
	if finder.Criteria != nil {
		for key, value := range finder.Criteria {
			values = append(values, value)
			if whereStatement == "" {
				whereStatement = fmt.Sprintf("%s = ?", key)
			} else {
				whereStatement = fmt.Sprintf("%s AND %s = ? ", whereStatement, key)
			}
		}
		whereStatement = " WHERE " + whereStatement
	}

	if finder.OrderBy != nil {
		for key, value := range finder.OrderBy {
			if orderByStatement == "" {
				orderByStatement = fmt.Sprintf("%s %s", key, value)
			} else {
				orderByStatement = fmt.Sprintf("%s , %s %s ", orderByStatement, key, value)
			}
		}
		orderByStatement = " ORDER BY " + orderByStatement
	}
	if finder.Limit != 0 {
		limitStatement = fmt.Sprintf(" LIMIT %d ", finder.Limit)
	}
	if finder.Offset != 0 {
		offsetStatement = fmt.Sprintf(" OFFSET %d ", finder.Offset)
	}
	query := []string{whereStatement, orderByStatement, limitStatement, offsetStatement}
	return strings.Join(query, ""), values
}

// UserRepository is a repository of users
type Repository struct {
	Connection *Connection
	IDField    string
	TableName  string
	Type       reflect.Type
	DM         *DataMapper
}

func NewRepository(Type reflect.Type, datamapper *DataMapper) *Repository {
	metadata, ok := datamapper.Metadatas[Type]
	if !ok {
		log.Fatalf("Datamapper cannot manage type %s", Type)
	}
	idField := metadata.FindIdColumn().Name
	if idField == "" {
		idField = metadata.FindIdColumn().StructField
	}
	return &Repository{datamapper.GetConnection(), idField, metadata.Table.Name, Type, datamapper}
}

// All finds all
func (repository *Repository) All(collection interface{}) error {
	return repository.Connection.Select(collection, fmt.Sprintf("SELECT %[1]s.* FROM %[1]s ", repository.TableName))
}

// Find finds a user by id
func (repository *Repository) Find(id interface{}, model interface{}) error {
	return repository.Connection.Get(model, fmt.Sprintf("SELECT %[1]s.* FROM %[1]s WHERE %[2]s =? ", repository.TableName, repository.IDField), id)
}

// FindBy find users by fields
func (repository *Repository) FindBy(finder Finder, collection interface{}) error {
	query, values := finder.AcceptRepository(repository)
	err := repository.Connection.Select(collection, fmt.Sprintf("SELECT * FROM %s %s;", repository.TableName, query), values...)
	if err != nil {
		return err
	}
	return nil
}

// DeleteAll deletes all models
func (repository *Repository) DeleteAll() error {
	result, err := repository.Connection.Exec(fmt.Sprintf("DELETE FROM %s;", repository.TableName))
	if err != nil {
		return err
	}
	if _, err := result.RowsAffected(); err != nil {
		return err
	}
	return nil
}

// Save saves a model
func (repository *Repository) Save(entity interface{}) error {
	if u, ok := interface{}(entity).(BeforeSaveCallback); ok == true {
		if err := u.BeforeSave(); err != nil {
			return err
		}
	}
	if repository.ResolveId(entity).(int64) == 0 {
		if u, ok := interface{}(entity).(BeforeCreateCallback); ok == true {
			if err := u.BeforeCreate(); err != nil {
				return err
			}
		}
		paths := []string{}
		values := []interface{}{}
		fieldMap := repository.DM.Metadatas[repository.Type].FieldMap(entity)

		for key, value := range fieldMap {
			if strings.ToLower(key) != strings.ToLower(repository.IDField) {
				paths = append(paths, key)
				values = append(values, value.Interface())
			}
		}
		query := fmt.Sprintf("INSERT INTO %s(%s) VALUES(%s);",
			repository.TableName,
			strings.Join(paths, ","),
			strings.Join(
				strings.Split(strings.Repeat("?", len(paths)), ""), ","))
		result, err := repository.Connection.Exec(query, values...)
		if err != nil {
			return err
		}
		id, err := result.LastInsertId()
		if err != nil {
			return err
		}
		return repository.Find(id, entity)
	}
	if u, ok := interface{}(entity).(BeforeUpdateCallback); ok == true {
		if err := u.BeforeUpdate(); err != nil {
			return err
		}
	}
	// See http://stackoverflow.com/questions/24318389/golang-elem-vs-indirect-in-the-reflect-package

	fieldMap := repository.DM.Metadatas[repository.Type].FieldMap(entity)
	attributes := map[string]interface{}{}
	for key, value := range fieldMap {
		if strings.ToLower(key) != strings.ToLower(repository.IDField) {
			attributes[key] = value.Interface()
		}
	}
	err := repository.UpdateAttribute(entity, attributes)
	if err != nil {
		return err
	}
	return nil
}

// UpdateAttribute update selected attributes
func (repository *Repository) UpdateAttribute(model interface{}, attributes map[string]interface{}) error {
	values := []interface{}{}
	setStatement := ""
	userType := reflect.TypeOf(model)
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
	id := repository.ResolveId(model)
	result, err := repository.Connection.Exec(fmt.Sprintf("UPDATE %s SET %s WHERE ID = ?;", repository.TableName, setStatement), append(values, id)...)
	if err != nil {
		return err
	}
	if _, err := result.RowsAffected(); err != nil {
		return err
	}
	err = repository.Find(id, model)
	if err != nil {
		return err
	}
	return nil
}

// Destroy removes a user
func (repository *Repository) Destroy(model interface{}) error {
	id := repository.ResolveId(model)
	result, err := repository.Connection.Exec(
		fmt.Sprintf("DELETE FROM %s WHERE %s.%s=?", repository.TableName, repository.IDField),
		id,
	)
	if err != nil {
		return err
	}
	if rows, err := result.RowsAffected(); err != nil {
		return err
	} else if rows <= 0 {
		return fmt.Errorf("User with ID %d could not be found and desftroyed.", id)
	}
	return nil
}

func (repository *Repository) ResolveId(model interface{}) interface{} {
	value := reflect.Indirect(reflect.ValueOf(model))
	return value.FieldByName(repository.IDField).Interface()
}

// CanManageEntity returns an error if the type of entity is not equal to
// repository.Type
func (repository *Repository) CanManageEntity(entity interface{}) error {

	if entityType := reflect.Indirect(reflect.ValueOf(entity)).Type(); entityType != repository.Type {
		return fmt.Errorf("Repository cannot manage entity of type %s", entityType)
	}
	return nil
}
