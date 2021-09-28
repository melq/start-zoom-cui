package repository

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/jmoiron/sqlx/types"
	"log"
	"strconv"
)

type Bit int

type Meet struct {
	Id		int				`db:"id"`
	Dispose types.BitBool 	`db:"dispose"`
	Name 	string			`db:"meet_name"`
	Url 	string			`db:"url"`
	Day 	sql.NullString 	`db:"day_of_week"`
	Date 	sql.NullString	`db:"meet_date"`
	STime 	string			`db:"start_time"`
	ETime	string			`db:"end_time"`
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
		" start_time time not null," +
		" end_time time not null)"

	tx := db.MustBegin()
	tx.MustExec("INSERT INTO users (name) VALUES (?)", user)
	tx.MustExec(schema)
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
			"(dispose, meet_name, url, meet_date, start_time, end_time)"+
			"VALUE (:dispose, :meet_name, :url, :meet_date, :start_time, :end_time)", meet)
	} else {
		_, err = db.NamedExec("INSERT INTO "+name+
			"(dispose, meet_name, url, day_of_week, start_time, end_time)"+
			"VALUE (:dispose, :meet_name, :url, :day_of_week, :start_time, :end_time)", meet)
	}
	if err != nil {
		log.Fatalln(err)
	}
}

func GetMeets(user string) []Meet {
	db := getDB()
	defer func(db *sqlx.DB) {
		err := db.Close()
		if err != nil {
			log.Fatalln()
		}
	}(db)
	
	var meetList []Meet
	err := db.Select(&meetList, "SELECT * FROM " + user + " ORDER BY dispose DESC, meet_date, day_of_week")
	if err != nil {
		log.Fatalln(err)
		return nil
	}
	return meetList
}

func GetMeetsWithOpts(user string, mode int) []Meet {
	db := getDB()
	defer func(db *sqlx.DB) {
		err := db.Close()
		if err != nil {
			log.Fatalln()
		}
	}(db)

	var meetList []Meet
	var err error
	if mode == 0 {
		err = db.Select(&meetList, "SELECT * FROM " + user + " WHERE dispose=true ORDER BY dispose DESC, meet_date, start_time")
	} else {
		err = db.Select(&meetList, "SELECT * FROM " + user + " WHERE dispose=false ORDER BY dispose DESC, day_of_week, start_time")
	}

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
	fmt.Println(meet)
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

	id := strconv.Itoa(meet.Id)
	tx := db.MustBegin()
	tx.MustExec("UPDATE " + user + " SET meet_name=? WHERE id='" + id + "'", meet.Name)
	tx.MustExec("UPDATE " + user + " SET url=? WHERE id='" + id + "'", meet.Url)
	if meet.Day.Valid == true {
		tx.MustExec("UPDATE " + user + " SET day_of_week=? WHERE id='" + id + "'", meet.Day.String)
		tx.MustExec("UPDATE " + user + " SET meet_date=? WHERE id='" + id + "'", nil)
	}
	if meet.Date.Valid == true {
		tx.MustExec("UPDATE " + user + " SET meet_date=? WHERE id='" + id + "'", meet.Date.String)
		tx.MustExec("UPDATE " + user + " SET day_of_week=? WHERE id='" + id + "'", nil)
	}
	tx.MustExec("UPDATE " + user + " SET start_time=? WHERE id='" + id + "'", meet.STime)
	tx.MustExec("UPDATE " + user + " SET end_time=? WHERE id='" + id + "'", meet.ETime)
	tx.MustExec("UPDATE " + user + " SET dispose=? WHERE id='" + id + "'", meet.Dispose)
	err := tx.Commit()
	if err != nil {
		log.Fatalln(err)
		return
	}
	fmt.Println(meet)
}

func DeleteMeet(user string, meetName string) {
	db := getDB()
	defer func(db *sqlx.DB) {
		err := db.Close()
		if err != nil {
			log.Fatalln()
		}
	}(db)

	_, err := db.Queryx("DELETE FROM " + user + " WHERE meet_name=? LIMIT 1", meetName)
	if err != nil {
		log.Fatalln(err)
	}
}