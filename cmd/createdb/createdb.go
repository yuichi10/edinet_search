package createdb

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/spf13/cobra"
)

/*
DB設計
update_dateテーブル。
今まで取得した日数のを記録しておく
将来的には基本的にはここの範囲外のデータのみを更新で追加していく。 forceオプションで強制的に取れるとかもしても良さそう。
最初は毎回データを削除して作成していく形式にする。
というか最初は全部消して作り直すならこのテーブルはいらないか。まずはこれは作らないでおく
- start_date: date
- end_date: date


documents_metaテーブル
実際にedinetから取得してきたドキュメントの情報を保管。ドキュメントと行っているが、有価証券の情報だけを保管する予定.
docIDをprimayにしておけば良さげ。
ordinanceCodeが010、form_codeが030000のデータのみをまずは取るようにしてみる。
- docID
- parentDocID
- filerName
- submitDateTime
- docDescription
*/

const EDINET_API_ENDPOINT = "https://api.edinet-fsa.go.jp/api/v2/documents.json"
const DB_NAME = "company.db"

// args
var strStartDate, strEndDate string

var startDate, endDate time.Time

type inputParams struct {
	startDate time.Time
	endDate   time.Time
}

func initDB() (*sql.DB, error) {
	os.Remove(DB_NAME)
	db, err := sql.Open("sqlite3", DB_NAME)
	if err != nil {
		return nil, fmt.Errorf("sqliteのファイルを参照するのに失敗しました。 %s", err.Error())
	}
	defer db.Close()

	sqlStmt := `
			create table documents (docID text not null primary key, secCode text, filerName text, docDescription text, submitDateTime text);
			`
	_, err = db.Exec(sqlStmt)
	if err != nil {
		return nil, fmt.Errorf("tableの作成に失敗しました。 %s", err.Error())
	}
	return db, err
}

func parseParams() (inputParams, error) {
	layout := "2006-01-02"
	var err error
	params := inputParams{}
	if params.startDate, err = time.Parse(layout, strStartDate); err != nil {
		return params, fmt.Errorf("最初の日付が不正な状態です")
	}

	if params.endDate, err = time.Parse(layout, strEndDate); err != nil {
		return params, fmt.Errorf("最後の日付が不正な状態です")
	}

	return params, nil
}

func NewCreateDBCmd() *cobra.Command {
	c := &cobra.Command{
		Use:   "createdb",
		Short: "書類の一覧を取得するDBを作成します。すでにDBがある場合は不要です。",
		Run: func(cmd *cobra.Command, args []string) {
			_, err := parseParams()
			if err != nil {
				log.Fatal(err.Error())
			}
			fmt.Println("空のDBを作成します。")
			_, err = initDB()
			if err != nil {
				log.Fatal(err.Error())
			}

		},
	}

	c.Flags().StringVarP(&strStartDate, "start-date", "s", "", "2023-12-12 の形式で追加してください")
	c.Flags().StringVarP(&strEndDate, "end-date", "e", "", "2023-12-12 の形式で追加してください")
	return c
}
