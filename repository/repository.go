package repository

import (
	"database/sql"
	"fmt"
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

func getDB() *sqlx.DB {
	db, err := sqlx.Open("mysql", "melq:pass@/meet")
	if err != nil {
		log.Fatalln(err)
	}
	return db
}

func CreateUser(user string) {
	db := getDB()
	defer func(db *sqlx.DB) {
		err := db.Close()
		if err != nil {
			log.Fatalln()
		}
	}(db)

	var schema = "CREATE TABLE IF NOT EXISTS " + user + // ユーザ毎のテーブルを作成するクエリ
		"(id int not null primary key auto_increment," +
		" dispose bit not null," +
		" meet_name varchar(32) not null," +
		" url varchar(256) not null," +
		" day_of_week enum('Sunday','Monday','Tuesday','Wednesday','Thursday','Friday','Saturday')," +
		" meet_date date," +
		" meet_time time)"

	tx := db.MustBegin()
	tx.MustExec("INSERT INTO users (name) VALUES (?)", user)
	log.Println(tx.MustExec(schema))
	err := tx.Commit()
	if err != nil {
		log.Fatalln(err)
	}
}

func MakeMeet(name string, meet Meet) {
	db := getDB()
	defer func(db *sqlx.DB) {
		err := db.Close()
		if err != nil {
			log.Fatalln()
		}
	}(db)

	var err error
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
	db := getDB()
	defer func(db *sqlx.DB) {
		err := db.Close()
		if err != nil {
			log.Fatalln()
		}
	}(db)
	
	var meetList []Meet
	err := db.Select(&meetList, "SELECT * FROM " + name + " ORDER BY dispose DESC, meet_date, day_of_week")
	if err != nil {
		log.Fatalln(err)
		return nil
	}
	return meetList
}

func GetMeet(user string, name string) Meet {
	db := getDB()
	defer func(db *sqlx.DB) {
		err := db.Close()
		if err != nil {
			log.Fatalln()
		}
	}(db)

	meet := Meet{}
	err := db.Get(&meet, "SELECT * FROM " + user + " WHERE meet_name=?", name)
	if err != nil {
		log.Fatalln(err)
	}
	return meet
}

func UpdateMeet(user string, meet Meet) {
	db := getDB()
	defer func(db *sqlx.DB) {
		err := db.Close()
		if err != nil {
			log.Fatalln()
		}
	}(db)

	fmt.Println(meet)
	tx := db.MustBegin()
	tx.MustExec("UPDATE " + user + " SET url=? WHERE meet_name='" + meet.Name + "'", meet.Url)
	if meet.Day.Valid == true {
		fmt.Println(meet.Day)
		tx.MustExec("UPDATE " + user + " SET day_of_week=? WHERE meet_name='" + meet.Name + "'", meet.Day.String)
	}
	if meet.Date.Valid == true {
		tx.MustExec("UPDATE " + user + " SET meet_date=? WHERE meet_name='" + meet.Name + "'", meet.Date.String)
	}
	tx.MustExec("UPDATE " + user + " SET meet_time=? WHERE meet_name='" + meet.Name + "'", meet.Time)
	err := tx.Commit()
	if err != nil {
		log.Fatalln(err)
		return
	}
	fmt.Println(meet)
}