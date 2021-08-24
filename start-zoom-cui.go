package main

import (
	"fmt"
	"github.com/jessevdk/go-flags"
	"os"
)

type Option struct {
	Start []bool `short:"s" long:"start" description:"近い会議を開始します"`
	Register []bool `short:"r" long:"register" description:"会議の予定を作成します"`
	List []bool `short:"l" long:"list" description:"登録されている会議の一覧を表示します"`
	Edit []bool `short:"e" long:"edit" description:"登録されている会議の編集・削除を行います"`
	Setting []bool `long:"setting" description:"設定を行います"`
}
var opts Option

func main() {
	_, err := flags.Parse(&opts)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if len(opts.Start) != 0 {
		fmt.Println("Start")
		// 会議開始機能
	} else if len(opts.Register) != 0 {
		fmt.Println("Register")
		// 会議登録機能
	} else if len(opts.List) != 0 {
		fmt.Println("List")
		// 登録会議閲覧機能
	} else if len(opts.Edit) != 0 {
		fmt.Println("Edit")
		// 登録会議編集・削除機能
	} else if len(opts.Setting) != 0 {
		fmt.Println("Setting")
		// 設定変更機能
	}
}