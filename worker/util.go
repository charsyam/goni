package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	influxlib "github.com/influxdata/influxdb/client/v2"
	"os"
)

func getDataSource() string {
	host := os.Getenv("GONI_MYSQL_HOST")
	port := os.Getenv("GONI_MYSQL_PORT")
	user := os.Getenv("GONI_MYSQL_USER")
	pass := os.Getenv("GONI_MYSQL_PASS")
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/goni_saas", user, pass, host, port)
}

func getInflux() (influxlib.Client, error) {
	influx, err := influxlib.NewHTTPClient(influxlib.HTTPConfig{
		Addr:     os.Getenv("GONI_INFLUX_HOST"),
		Username: os.Getenv("GONI_INFLUX_USER"),
		Password: os.Getenv("GONI_INFLUX_PASS"),
	})
	if err != nil {
		return nil, err
	}
	return influx, nil
}

func getMySQL() (*sql.DB, error) {
	db, err := sql.Open("mysql", getDataSource())
	if err != nil {
		return nil, err
	}
	return db, nil
}
