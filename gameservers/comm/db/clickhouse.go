package db

import (
	"database/sql"
	"fmt"
	_ "github.com/ClickHouse/clickhouse-go"
	"log"
)

var (
	clickhouseDB *sql.DB
)

func DialToClickHouse(clickHouseAddr string) error {
	conn, err := sql.Open("clickhouse", clickHouseAddr)
	if err != nil {
		return fmt.Errorf("failed to connect to ClickHouse: %w", err)
	}

	if err := conn.Ping(); err != nil {
		conn.Close()
		return fmt.Errorf("failed to ping ClickHouse: %w", err)
	}
	clickhouseDB = conn

	return nil
}

func ClickHouseCollection(dbName string) (*sql.DB, error) {
	if dbName != "" {
		_, err := clickhouseDB.Exec(fmt.Sprintf("USE %s;", dbName))
		if err != nil {
			log.Fatalf("查询失败: %v", err)
		}
		return clickhouseDB, err
	}
	return clickhouseDB, nil
}

func ClickHouseClose() {
	clickhouseDB.Close()
}
