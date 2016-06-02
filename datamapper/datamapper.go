package datamapper

import (
	"encoding/json"
	"fmt"
	"reflect"
)

type Table struct {
	Name string
}
type Column struct {
	Id          bool
	StructField string
	Name        string
}

// DataMapperMetadata represent metadatas for a DB Table
type DataMapperMetadata struct {
	Entity  string
	Table   Table
	Columns []Column
}

// DataMapperFromMetadataFrom creates a datamapper metadata from a json string
// or returns an error
func DataMapperMetadataFrom(jsonString string) (DataMapperMetadata, error) {
	var dmm DataMapperMetadata
	err := json.Unmarshal([]byte(jsonString), &dmm)
	return dmm, err
}

func (dmm DataMapperMetadata) FindIdColumn() Column {
	var column Column
	for _, value := range dmm.Columns {
		if value.Id {
			column = value
			break
		}
	}
	return column
}

func (dmm DataMapperMetadata) FieldMap(entity interface{}) (fieldMap map[string]reflect.Value) {
	value := reflect.Indirect(reflect.ValueOf(entity))
	fieldMap = map[string]reflect.Value{}
	for _, column := range dmm.Columns {
		name := column.Name
		if name == "" {
			name = column.StructField
		}
		fieldMap[name] = value.FieldByName(column.StructField)
	}
	return
}

type DataMapper struct {
	Connection *Connection
	Metadatas  map[reflect.Type]DataMapperMetadata
}

func NewDataMapper(connection *Connection) *DataMapper {
	return &DataMapper{connection, map[reflect.Type]DataMapperMetadata{}}
}

func (dm DataMapper) Register(entity interface{}) error {
	if e, ok := entity.(DataMapperMetadataProvider); ok {
		dm.Metadatas[reflect.TypeOf(entity)] = e.DataMapperMetaData()
		return nil
	}
	return fmt.Errorf("Cannot create metadata from Entity %v .", entity)
}

func (dm *DataMapper) GetRepository(entity interface{}) (*Repository, error) {
	Type := reflect.TypeOf(entity)
	if _, ok := dm.Metadatas[Type]; ok {
		return NewRepository(Type, dm), nil
	}
	return nil, fmt.Errorf("Metadata not found for type %s .", Type)
}

func (dm *DataMapper) GetConnection() *Connection {
	return dm.Connection
}
