package repository

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/jmoiron/sqlx/types"
	"log"
	"strconv"
)

var DayOfWeekString = [7]string{
	"Sunday",
	"Monday",
	"Tuesday",
	"Wednesday",
	"Thursday",
	"Friday",
	"Saturday",
}

type Meet struct {
	Id     int            `db:"id"`
	Weekly types.BitBool  `db:"weekly"`
	Name   string         `db:"meet_name"`
	Url    string         `db:"url"`
	Day    sql.NullString `db:"day_of_week"`
	Date   sql.NullString `db:"meet_date"`
	STime  string         `db:"start_time"`
	ETime  string         `db:"end_time"`
}

func getDB() *sqlx.DB {
	db, err := sqlx.Open("mysql", "melq:pass@/meet")
	if err != nil {
		log.Fatalln(err)
	}
	return db
}

func CreateUser(user string) int {
	db := getDB()
	defer func(db *sqlx.DB) {
		err := db.Close()
		if err != nil {
			log.Fatalln("getDB:", err)
		}
	}(db)

	type UserData struct {
		Id   int    `db:"id"`
		Name string `db:"name"`
	}

	var users []UserData
	err := db.Select(&users, "SELECT * FROM users WHERE name=(?)", user)
	if err != nil {
		log.Fatalln("CreateUser/getUsers:", err)
	}
	for _, v := range users {
		if v.Name == user {
			return 0
		}
	}

	var schema = "CREATE TABLE IF NOT EXISTS " + user + // ユーザ毎のテーブルを作成するクエリ
		"(id int not null primary key auto_increment," +
		" weekly bit not null," +
		" meet_name varchar(32) not null," +
		" url varchar(256) not null," +
		" day_of_week enum('Sunday','Monday','Tuesday','Wednesday','Thursday','Friday','Saturday')," +
		" meet_date date," +
		" start_time time not null," +
		" end_time time not null)"

	tx := db.MustBegin()
	tx.MustExec("INSERT INTO users (name) VALUES (?)", user)
	tx.MustExec(schema)
	err = tx.Commit()
	if err != nil {
		log.Fatalln("CreateUser:", err)
	}
	return 1
}

func MakeMeet(name string, meet Meet) {
	db := getDB()
	defer func(db *sqlx.DB) {
		err := db.Close()
		if err != nil {
			log.Fatalln(err)
		}
	}(db)

	var err error
	if meet.Weekly {
		_, err = db.NamedExec("INSERT INTO "+name+
			"(weekly, meet_name, url, day_of_week, start_time, end_time)"+
			"VALUE (:weekly, :meet_name, :url, :day_of_week, :start_time, :end_time)", meet)
	} else {
		_, err = db.NamedExec("INSERT INTO "+name+
			"(weekly, meet_name, url, meet_date, start_time, end_time)"+
			"VALUE (:weekly, :meet_name, :url, :meet_date, :start_time, :end_time)", meet)
	}
	if err != nil {
		log.Fatalln("MakeMeet:", err)
	}
}

func GetMeets(user string) []Meet {
	db := getDB()
	defer func(db *sqlx.DB) {
		err := db.Close()
		if err != nil {
			log.Fatalln(err)
		}
	}(db)

	var meetList []Meet
	err := db.Select(&meetList, "SELECT * FROM "+user+" ORDER BY weekly DESC, meet_date, day_of_week")
	if err != nil {
		log.Fatalln("GetMeets:", err)
		return nil
	}
	return meetList
}

func GetMeetsWithOpts(user string, mode int) []Meet {
	db := getDB()
	defer func(db *sqlx.DB) {
		err := db.Close()
		if err != nil {
			log.Fatalln(err)
		}
	}(db)

	var meetList []Meet
	var err error
	if mode == 0 {
		err = db.Select(&meetList, "SELECT * FROM "+user+" WHERE weekly=true ORDER BY weekly DESC, meet_date, start_time")
	} else {
		err = db.Select(&meetList, "SELECT * FROM "+user+" WHERE weekly=false ORDER BY weekly DESC, day_of_week, start_time")
	}

	if err != nil {
		log.Fatalln("GetMeetsWithOpts:", err)
		return nil
	}
	return meetList
}

func GetMeet(user string, name string) Meet {
	db := getDB()
	defer func(db *sqlx.DB) {
		err := db.Close()
		if err != nil {
			log.Fatalln(err)
		}
	}(db)

	meet := Meet{}
	err := db.Get(&meet, "SELECT * FROM "+user+" WHERE meet_name=?", name)
	if err != nil {
		log.Fatalln("GetMeet:", err)
	}
	return meet
}

func UpdateMeet(user string, meet Meet) {
	db := getDB()
	defer func(db *sqlx.DB) {
		err := db.Close()
		if err != nil {
			log.Fatalln(err)
		}
	}(db)

	id := strconv.Itoa(meet.Id)
	tx := db.MustBegin()
	tx.MustExec("UPDATE "+user+" SET meet_name=? WHERE id='"+id+"'", meet.Name)
	tx.MustExec("UPDATE "+user+" SET url=? WHERE id='"+id+"'", meet.Url)
	if meet.Day.Valid == true {
		tx.MustExec("UPDATE "+user+" SET day_of_week=? WHERE id='"+id+"'", meet.Day.String)
		tx.MustExec("UPDATE "+user+" SET meet_date=? WHERE id='"+id+"'", nil)
	}
	if meet.Date.Valid == true {
		tx.MustExec("UPDATE "+user+" SET meet_date=? WHERE id='"+id+"'", meet.Date.String)
		tx.MustExec("UPDATE "+user+" SET day_of_week=? WHERE id='"+id+"'", nil)
	}
	tx.MustExec("UPDATE "+user+" SET start_time=? WHERE id='"+id+"'", meet.STime)
	tx.MustExec("UPDATE "+user+" SET end_time=? WHERE id='"+id+"'", meet.ETime)
	tx.MustExec("UPDATE "+user+" SET weekly=? WHERE id='"+id+"'", meet.Weekly)
	err := tx.Commit()
	if err != nil {
		log.Fatalln("UpdateMeet:", err)
		return
	}
}

func DeleteMeet(user string, meetName string) {
	db := getDB()
	defer func(db *sqlx.DB) {
		err := db.Close()
		if err != nil {
			log.Fatalln(err)
		}
	}(db)

	_, err := db.Queryx("DELETE FROM "+user+" WHERE meet_name=? LIMIT 1", meetName)
	if err != nil {
		log.Fatalln("DeleteMeet:", err)
	}
}
