package db

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"time"
)

const (
	ROLE_OWNER = "OWNER"
	ROLE_GUEST = "GUEST"
)

var dbase *sql.DB
var tablesCreations = map[string]string{
	"auth": "create table auth (chat_id integer primary key, role varchar(255) not null, valid_until datetime not null)",
}
var validDurations = map[string]time.Duration{
	ROLE_OWNER: 90 * 24 * time.Hour,
	ROLE_GUEST: time.Hour,
}

func Init() error {
	db, err := sql.Open("sqlite3", "raspi-manage-bot.db")
	if err != nil {
		return err
	}

	dbase = db

	if err := migrate(); err != nil {
		return err
	}

	return nil
}

func Close() error {
	return dbase.Close()
}

func migrate() error {
	if err := dbase.Ping(); err != nil {
		return err
	}

	for table := range tablesCreations {
		row := dbase.QueryRow("select name from sqlite_master where type = 'table' and name = $1", table)
		var tableName string
		if err := row.Scan(&tableName); err != nil {
			_, err = dbase.Exec(tablesCreations[table])
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func InsertAuth(chatID int64, role string) error {
	_, err := dbase.Exec("insert or replace into auth (chat_id, role, valid_until) values ($1, $2, $3)", chatID, role, time.Now().Add(validDurations[role]))
	if err != nil {
		return err
	}

	return nil
}

func GetRole(chatID int64) string {
	row := dbase.QueryRow("select role from auth where chat_id = $1 and valid_until >= $2", chatID, time.Now())
	var role string
	err := row.Scan(&role)
	if err != nil {
		return ""
	}

	return role
}
