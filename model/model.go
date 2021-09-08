package model

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"log"
)

func CreateUser(name string) {
	db, err := sqlx.Open("mysql", "melq:pass@/meet")
	if err != nil {
		log.Fatalln(err)
	}
	defer func(db *sqlx.DB) {
		err := db.Close()
		if err != nil {
			log.Fatalln(err)
		}
	}(db)

	//var schema = "CREATE TABLE " + name + "" // ユーザ毎のテーブルを作成するクエリ

	tx := db.MustBegin()
	tx.MustExec("INSERT INTO users (name) VALUES (?)", name)
	err = tx.Commit()
	if err != nil {
		log.Fatalln(err)
	}
}
