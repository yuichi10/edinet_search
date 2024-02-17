package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

/*
DB設計
updateテーブル
今まで取得した日数のを記録しておく
将来的には基本的にはここの範囲外のデータのみを更新で追加していく。 forceオプションで強制的に取れるとかもしても良さそう。
まずは全ての範囲のデータを取得して更新していく。
- start_date: date
- end_date: date


documents_metaテーブル
実際にedinetから取得してきたドキュメントの情報を保管。ドキュメントと行っているが、有価証券の情報だけを保管する予定.
docIDをprimayにしておけば良さげ。
- docID
- parentDocID
- filerName
- submitDateTime
- docDescription
*/

func NewCreateDBCmd() *cobra.Command {
	c := &cobra.Command{
		Use: "createdb",
		Short: "書類の一覧を取得するDBを作成します。すでにDBがある場合は不要です。",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("create db")
			fmt.Println(viper.GetString("api.token"))
		},
	}
	return c
}
