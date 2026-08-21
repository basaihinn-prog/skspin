package sqlite

import (
	"database/sql"
	"fmt"
	"game/duck/logger"
	"log"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type DB struct {
	*sql.DB
	dbfile string
}

func NewDB(dbfile string) *DB {
	db, err := sql.Open("sqlite3", dbfile)
	if err != nil {
		log.Fatal(err)
	}
	tmpDB := &DB{
		DB:     db,
		dbfile: dbfile,
	}
	go tmpDB.keeplive()
	return tmpDB
}

func (this *DB) keeplive() {
	for {
		err := this.Ping()
		if err != nil {
			client, err := sql.Open("sqlite3", this.dbfile)
			if err != nil {
				logger.Err(err)
			} else {
				this.DB = client
			}
		}
		time.Sleep(1 * time.Second)
	}
}

func (this *DB) CreateOriginTable(tableName string) error {
	sql := fmt.Sprintf("create table if not exists %v (id integer primary key, md5s text unique,pan text not null);", tableName)
	_, err := this.Exec(sql)
	return err
}

func (this *DB) CreateReleaseTable(tableName string) error {
	sql := fmt.Sprintf("create table if not exists %v (id integer primary key,count integer,pan text not null);", tableName)
	_, err := this.Exec(sql)
	return err
}

// 需要预先编译
func (this *DB) InsertOriginTableSql(tableName string) string {
	return fmt.Sprintf("insert into %v(md5s,pan)values(?,?);", tableName)
}

// 需要预先编译
func (this *DB) InsertReleaseTableSql(tableName string) string {
	return fmt.Sprintf("insert into %v(count,pan)values(?,?);", tableName)
}
func (this *DB) TruncateTable(tableName string) error {
	// sql := fmt.Sprintf("truncate %v;", tableName)
	sql := fmt.Sprintf("delete from %v;", tableName)
	_, err := this.Exec(sql)
	return err
}

func (this *DB) GetTableCount(tablename string) (int, error) {
	var count int
	err := this.QueryRow(fmt.Sprintf("select count(1) from %v;", tablename)).Scan(&count)
	return count, err
}

func (this *DB) GetDataFromTableByOffset(tablename string, start, offset int) (result []string, err error) {
	rows, err := this.Query(fmt.Sprintf("select pan from %v limit %v,%v;", tablename, start, offset))
	result = ResultPan(rows)
	return
}

func (this *DB) GetIDsFromTableByOffset(tablename string, start, offset int) (result []int, err error) {
	rows, err := this.Query(fmt.Sprintf("select id from %v limit %v,%v;", tablename, start, offset))
	result = ResultID(rows)
	return
}

func (this *DB) GetIDSFromTable(tablename string) (result []int, err error) {
	var start int
	const step = 10000
	for {
		results, err1 := this.GetIDsFromTableByOffset(tablename, start, step)
		if err1 != nil {
			logger.Fatal(err1)
			err = err1
			return
		}
		length := len(results)
		if length > 0 {
			result = append(result, results...)
		} else {
			break
		}
		start += length
	}
	return
}

func (this *DB) DeleteFromTableByIDS(tablename string, ids ...int) error {
	tx, err := this.Begin()
	if err != nil {
		return err
	}
	stmt, err := tx.Prepare(fmt.Sprintf("delete from %v where id = ?", tablename))
	if err != nil {
		tx.Rollback()
		return err
	}
	for _, id := range ids {
		_, err := stmt.Exec(id)
		if err != nil {
			tx.Rollback()
			return err
		}
	}
	tx.Commit()
	return nil
}

func (this *DB) GetDataFromTableByID(tablename string, id int) (result string, err error) {
	err = this.QueryRow(fmt.Sprintf("select pan from %v where id = %v;", tablename, id)).Scan(&result)
	return
}

func ResultPan(query *sql.Rows) []string {
	results := make([]string, 0)
	for query.Next() { //循环，让游标往下移动
		var pan string
		if err := query.Scan(&pan); err != nil {
			logger.Err(err)
			return results
		}
		results = append(results, pan)
	}
	return results
}

func ResultID(query *sql.Rows) []int {
	results := make([]int, 0)
	for query.Next() { //循环，让游标往下移动
		var ID int
		if err := query.Scan(&ID); err != nil {
			logger.Err(err)
			return results
		}
		results = append(results, ID)
	}
	return results
}
