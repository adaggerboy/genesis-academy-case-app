package postgres

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/adaggerboy/genesis-academy-case-app/models/config"
	"github.com/adaggerboy/genesis-academy-case-app/pkg/database"

	_ "github.com/lib/pq"
)

type PostgresConnection struct {
	database.IDatabase
	db      *sql.DB
	timeout time.Duration

	statements map[string]*sql.Stmt
}

func init() {
	database.GetDatabaseFabric().RegisterDatabaseFabric("postgres", getDI)
}

func createConnect(endpoint config.DatabaseEndpointConfig) (conn *sql.DB, err error) {
	conn, err = sql.Open("postgres", fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
		endpoint.User,
		endpoint.Password,
		endpoint.Host,
		endpoint.Port,
		endpoint.Database))
	return
}

func getDI(endpoint config.DatabaseEndpointConfig) (database.IDatabase, error) {
	conn, err := createConnect(endpoint)
	if err != nil {
		return nil, fmt.Errorf("postgres database connect error: %s", err)
	}
	chk, err := conn.Query("select 1")
	if err != nil {
		return nil, fmt.Errorf("postgres database connect error: %s", err)
	}
	if chk == nil || !chk.Next() {
		return nil, fmt.Errorf("postgres database connect error: some error idk")
	}
	chk.Close()

	reader := &PostgresConnection{db: conn,
		timeout: time.Duration(endpoint.Timeout)}
	reader.deploySchemas()
	err = reader.prepareStatements()
	if err != nil {
		return nil, fmt.Errorf("postgres database prepare error: %s", err)
	}
	return reader, nil
}
