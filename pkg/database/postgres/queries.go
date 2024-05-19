package postgres

import (
	"database/sql"
	_ "embed"
	"fmt"
)

//go:embed schema.sql
var schemaCreateQuery string

func getPreparedMap() map[string]string {
	return map[string]string{
		"subscribe":         `insert into subscriptions (email) values ($1) on conflict do nothing returning 1`,
		"unsubscribe":       `delete from subscriptions where email = $1`,
		"get_subscriptions": `select * from subscriptions`,
	}
}

func (db *PostgresConnection) deploySchemas() error {
	_, err := db.db.Exec(schemaCreateQuery)
	return err
}

func (db *PostgresConnection) prepareStatements() error {
	db.statements = map[string]*sql.Stmt{}

	for k, v := range getPreparedMap() {
		stmt, err := db.db.Prepare(v)
		if err != nil {
			return fmt.Errorf("prepare statements: %s", err)
		}
		db.statements[k] = stmt
	}
	return nil
}

func (db *PostgresConnection) CreateSubscription(email string) (created bool, err error) {
	rows, err := db.statements["subscribe"].Query(email)
	if err != nil {
		return
	}
	defer rows.Close()
	created = rows.Next()
	return
}

func (db *PostgresConnection) DeleteSubscription(email string) (err error) {
	_, err = db.statements["unsubscribe"].Exec(email)
	return
}

func (db *PostgresConnection) GetSubscriptions() (s []string, err error) {
	s = []string{}
	rows, err := db.statements["get_subscriptions"].Query()
	if err != nil {
		return
	}
	defer rows.Close()
	for rows.Next() {
		var tmp string
		err = rows.Scan(&tmp)
		if err != nil {
			return
		}
		s = append(s, tmp)
	}
	return
}
