package main

import (
	"database/sql"
	"fmt"
	"github.com/jessevdk/go-flags"
	"os"
	"start-zoom-cui/repository"
	"strconv"
)

type Option struct {
	Register bool `short:"r" long:"register" description:"アカウントを作成します"`
	Start bool `short:"s" long:"start" description:"近い会議を開始します"`
	Make bool `short:"m" long:"make" description:"会議の予定を作成します"`
	List bool `short:"l" long:"list" description:"登録されている会議の一覧を表示します"`
	Edit bool `short:"e" long:"edit" description:"登録されている会議の編集・削除を行います"`
	Setting bool `long:"setting" description:"設定を行います"`
	User string `short:"u" description:"ユーザ名を入力します"`
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
		makeMeets()
		// 会議登録機能 // 情報入力の仕様要検討
	} else if opts.Start {
		fmt.Println("Start", opts.User)
		// 会議開始機能
	} else if opts.List {
		fmt.Println("List", opts.User)
		showList()
		// 登録会議閲覧機能
	} else if opts.Edit {
		fmt.Println("Edit", opts.User)
		// 登録会議編集・削除機能
	} else if opts.Setting {
		fmt.Println("Setting", opts.User)
		// 設定変更機能
	} else {
		flags.NewParser(&opts, flags.Default).WriteHelp(os.Stdout)
		os.Exit(0)
	}
}

func makeMeets() { // TEST
	for i := 0; i < 5; i++ {
		repository.MakeMeet(opts.User, repository.Meet{
			Dispose: false,
			Name:    "test" + strconv.Itoa(i),
			Url:     "example.com" + "/" + strconv.Itoa(i),
			Day:     sql.NullString{String: "Sunday", Valid: true},
			Date:    sql.NullString{Valid: false},
			Time:    fmt.Sprintf("%06d", 150500 + i),
		})

		repository.MakeMeet(opts.User, repository.Meet{
			Dispose: true,
			Name:    "test" + strconv.Itoa(-i),
			Url:     "example.com" + "/" + strconv.Itoa(-i),
			Day:     sql.NullString{Valid: false},
			Date:    sql.NullString{String: strconv.Itoa(210916 - i), Valid: true},
			Time:    fmt.Sprintf("%06d", 30500 + i),
		})
	}
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
		fmt.Println(" " + meet.Time + "\n")
	}
}