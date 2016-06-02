package datamapper

import (
	"database/sql"

	"github.com/jmoiron/sqlx"
)

type Logger interface {
	Log(arguments ...interface{})
}

type ConnectionOptions struct {
	Logger
}
type Connection struct {
	db         *sqlx.DB
	driverName string
	Options    *ConnectionOptions
}

func NewConnection(driverName string, DB *sql.DB) *Connection {
	return &Connection{sqlx.NewDb(DB, driverName), driverName, &ConnectionOptions{}}
}

func NewConnectionWithOptions(driverName string, DB *sql.DB, options *ConnectionOptions) *Connection {
	return &Connection{sqlx.NewDb(DB, driverName), driverName, options}
}

func (connection *Connection) DriverName() string {
	return connection.driverName
}

func (connection *Connection) Exec(query string, parameters ...interface{}) (sql.Result, error) {
	defer connection.Options.Logger.Log(append([]interface{}{query}, parameters...))
	return connection.db.Unsafe().Exec(query, parameters...)
}

func (connection *Connection) Select(records interface{}, query string, parameters ...interface{}) error {
	defer connection.Options.Logger.Log(append([]interface{}{query}, parameters...))
	return connection.db.Unsafe().Select(records, query, parameters...)

}

func (connection *Connection) Get(record interface{}, query string, parameters ...interface{}) error {
	defer connection.Options.Logger.Log(append([]interface{}{query}, parameters...))
	return connection.db.Unsafe().Get(record, query, parameters...)
}

func (connection *Connection) Log(messages ...interface{}) {
	if connection.Options.Logger != nil {
		connection.Options.Logger.Log(messages...)
	}
}
