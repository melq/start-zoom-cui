package repository

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"log"
)

type DayOfWeek int
const (
	Sunday = iota
	Monday
	Tuesday
	Wednesday
	Thursday
	Friday
	Saturday
)

type Meet struct {
	Dispose bool 		`db:"dispose"`
	Name 	string		`db:"meet_name"`
	Url 	string		`db:"url"`
	Day 	DayOfWeek 	`db:"day_of_week"`
	Date 	string		`db:"meet_date"`
	Time 	string		`db:"meet_time"`
}

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
	} (db)

	var schema = "CREATE TABLE IF NOT EXISTS " + name + // ユーザ毎のテーブルを作成するクエリ
		"(id int not null primary key auto_increment," +
		" dispose bit not null," +
		" meet_name varchar(32) not null," +
		" url varchar(256) not null," +
		" day_of_week enum('Sunday','Monday','Tuesday','Wednesday','Thursday','Friday','Saturday')," +
		" meet_date date," +
		" meet_time time)"

	tx := db.MustBegin()
	tx.MustExec("INSERT INTO users (name) VALUES (?)", name)
	log.Println(tx.MustExec(schema))
	err = tx.Commit()
	if err != nil {
		log.Fatalln(err)
	}
}

