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

	var schema = "CREATE TABLE " + name + // ユーザ毎のテーブルを作成するクエリ
		"(id int not null primary key auto_increment," +
		" meet_name varchar(32)," +
		" url varchar(256))"

	tx := db.MustBegin()
	tx.MustExec("INSERT INTO users (name) VALUES (?)", name)
	tx.MustExec(schema)
	err = tx.Commit()
	if err != nil {
		log.Fatalln(err)
	}
}
