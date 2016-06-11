package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"os"
)

func getDataSource() string {
	host := os.Getenv("GONI_MYSQL_HOST")
	port := os.Getenv("GONI_MYSQL_PORT")
	user := os.Getenv("GONI_MYSQL_USER")
	pass := os.Getenv("GONI_MYSQL_PASS")
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/goni_saas", user, pass, host, port)
}

func getMySQL() (*sql.DB, error) {
	db, err := sql.Open("mysql", getDataSource())
	if err != nil {
		return nil, err
	}
	return db, nil
}
