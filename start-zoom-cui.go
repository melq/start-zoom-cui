package main

import (
	"fmt"
	"github.com/jessevdk/go-flags"
	"os"
	"start-zoom-cui/repository"
)

type Option struct {
	Start bool `short:"s" long:"start" description:"近い会議を開始します"`
	Register bool `short:"r" long:"register" description:"会議の予定を作成します"`
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

	if opts.Start {
		fmt.Println("Start", opts.User)
		// 会議開始機能
	} else if opts.Register {
		fmt.Println("Register", opts.User)
		// 会議登録機能
	} else if opts.List {
		fmt.Println("List", opts.User)
		// 登録会議閲覧機能
	} else if opts.Edit {
		fmt.Println("Edit", opts.User)
		// 登録会議編集・削除機能
	} else if opts.Setting {
		fmt.Println("Setting", opts.User)
		// 設定変更機能
	} else {
		repository.CreateUser("test")
		flags.NewParser(&opts, flags.Default).WriteHelp(os.Stdout)
		os.Exit(0)
	}
}