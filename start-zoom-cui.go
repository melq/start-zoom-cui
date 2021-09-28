package main

import (
	"database/sql"
	"fmt"
	"github.com/jessevdk/go-flags"
	"github.com/jmoiron/sqlx/types"
	"os"
	"start-zoom-cui/repository"
)

type Option struct {
	Register bool	`short:"r" long:"register" description:"アカウントを作成します"`
	Start bool		`short:"s" long:"start" description:"近い会議を開始します"`
	Make bool		`short:"m" long:"make" description:"会議の予定を作成します"`
	List bool		`short:"l" long:"list" description:"登録されている会議の一覧を表示します"`
	Edit bool		`short:"e" long:"edit" description:"登録されている会議の編集を行います"`
	Delete bool		`short:"d" long:"delete" description:"登録されている会議の削除を行います"`
	Setting bool	`long:"setting" description:"設定を行います"`
	User string		`short:"u" description:"ユーザ名を入力します"`

	Dispose bool	`long:"dispose" description:"一度のみの会議の場合指定します"`
	Name string		`long:"name" description:"会議の名前を入力します"`
	Url string		`long:"url" description:"会議のURLを入力します"`
	Day string		`long:"day" description:"定期的な会議の場合、その曜日を入力します(形式: Sunday, Monday..)"`
	Date string		`long:"date" description:"定期的でない単発の会議の場合、その日付を入力します(形式: 2021年9月20日 -> 210920)"`
	STime string	`long:"stime" description:"会議の開始時刻を入力します(形式: 15:00:00 -> 150000)"`
	ETime string	`long:"etime" description:"会議の開始時刻を入力します(形式: 15:00:00 -> 150000)"`
}
var opts Option

func main() {
	_, err := flags.Parse(&opts)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if opts.Register {
		fmt.Println("Register", opts.User)
		repository.CreateUser(opts.User)

	} else if opts.Make {
		fmt.Println("Make", opts.User)
		meet := repository.Meet{
			Dispose: types.BitBool(opts.Dispose),
			Name: opts.Name,
			Url: opts.Url,
			STime: opts.STime,
			ETime: opts.ETime,
		}
		if meet.Dispose {
			meet.Date = sql.NullString{String: opts.Date, Valid: true}
		} else {
			meet.Day = sql.NullString{String: opts.Day, Valid: true}
		}
		makeMeet(opts, meet) // 情報入力の仕様要検討

	} else if opts.List {
		fmt.Println("List", opts.User)
		showList()

	} else if opts.Start {
		fmt.Println("Start", opts.User)
		startMeet(opts)
		// 会議開始機能

	} else if opts.Edit {
		fmt.Println("Edit", opts.User)
		editMeet(opts)

	} else if opts.Delete {
		fmt.Println("Delete", opts.User)
		deleteMeet(opts)

	} else if opts.Setting {
		fmt.Println("Setting", opts.User)
		// 設定変更機能

	} else {
		flags.NewParser(&opts, flags.Default).WriteHelp(os.Stdout)
		os.Exit(0)
	}
}

func makeMeet(opts Option, meet repository.Meet) {
	repository.MakeMeet(opts.User, meet)

	/*for i := 0; i < 3; i++ { // TEST
		repository.MakeMeet(opts.User, repository.Meet{
			Dispose:	false,
			Name:    	"test" + strconv.Itoa(i),
			Url:     	"example.com" + "/" + strconv.Itoa(i),
			Day:     	sql.NullString{String: "Sunday", Valid: true},
			Date:    	sql.NullString{Valid: false},
			STime:   	fmt.Sprintf("%06d", 150500 + i),
			ETime:		fmt.Sprintf("%06d", 173000 + i),
		})

		repository.MakeMeet(opts.User, repository.Meet{
			Dispose:	true,
			Name:		"test" + strconv.Itoa(-i),
			Url:		"example.com" + "/" + strconv.Itoa(-i),
			Day:		sql.NullString{Valid: false},
			Date:		sql.NullString{String: strconv.Itoa(210916 - i), Valid: true},
			STime:		fmt.Sprintf("%06d", 30500 + i),
			ETime: 		fmt.Sprintf("%06d", 53000 + i),
		})
	}*/
}

func showList() {
	meetList := repository.GetMeets(opts.User)
	for _, meet := range meetList {
		fmt.Println("会議名: " + meet.Name)
		fmt.Println("URL: " + meet.Url)
		if meet.Dispose {
			fmt.Print("日時: " + meet.Date.String)
		} else {
			fmt.Print("曜日: " + meet.Day.String)
		}
		fmt.Print(" " + meet.STime + " - ")
		fmt.Println(meet.ETime + "\n")
	}
}

func startMeet(opts Option) {
	/*meetList := repository.GetMeetsWithOpts(opts.User, 0)
	now := time.Now()
	var nextMeet repository.Meet
	for i, meet := range meetList {

	}

	meetList = repository.GetMeetsWithOpts(opts.User, 1)*/
}

func editMeet(opts Option) {
	meet := repository.GetMeet(opts.User, opts.Name)
	if opts.Url != "" { meet.Url = opts.Url }
	if len(opts.Day) > 0 {
		meet.Day = sql.NullString{ String: opts.Day, Valid: true }
		meet.Dispose = false
		meet.Date.Valid = false
	}
	if len(opts.Date) > 0 { // 曜日よりも日付指定を優先するので、こちらが後
		meet.Date = sql.NullString{ String: opts.Date, Valid: true }
		meet.Dispose = true
		meet.Day.Valid = false
	}
	if opts.STime != "" { meet.STime = opts.STime }
	if opts.ETime != "" { meet.ETime = opts.ETime }
	repository.UpdateMeet(opts.User, meet)
}

func deleteMeet(opts Option) {
	repository.DeleteMeet(opts.User, opts.Name)
}