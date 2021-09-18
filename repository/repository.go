package repository

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/jmoiron/sqlx/types"
	"log"
)

type Bit int

type Meet struct {
	Id		int				`db:"id"`
	Dispose types.BitBool 	`db:"dispose"`
	Name 	string			`db:"meet_name"`
	Url 	string			`db:"url"`
	Day 	sql.NullString 	`db:"day_of_week"`
	Date 	sql.NullString	`db:"meet_date"`
	Time 	string			`db:"meet_time"`
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

func MakeMeet(name string, meet Meet) {
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

	if meet.Dispose {
		_, err = db.NamedExec("INSERT INTO "+name+
			"(dispose, meet_name, url, meet_date, meet_time)"+
			"VALUE (:dispose, :meet_name, :url, :meet_date, :meet_time)", meet)
	} else {
		_, err = db.NamedExec("INSERT INTO "+name+
			"(dispose, meet_name, url, day_of_week, meet_time)"+
			"VALUE (:dispose, :meet_name, :url, :day_of_week, :meet_time)", meet)
	}
	if err != nil {
		log.Fatalln(err)
	}
}

func GetMeets(name string) []Meet {
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
	
	var meetList []Meet
	err = db.Select(&meetList, "SELECT * FROM "+name+" ORDER BY dispose DESC")
	if err != nil {
		log.Fatalln(err)
		return nil
	}
	return meetList
}