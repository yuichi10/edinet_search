package createdb

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
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

type Documents struct {
	Metadata struct {
		Title     string `json:"title"`
		Parameter struct {
			Date string `json:"date"`
			Type string `json:"type"`
		} `json:"parameter"`
		Resultset struct {
			Count int `json:"count"`
		} `json:"resultset"`
		ProcessDateTime string `json:"processDateTime"`
		Status          string `json:"status"`
		Message         string `json:"message"`
	} `json:"metadata"`
	Results []struct {
		SeqNumber            int         `json:"seqNumber"`
		DocID                string      `json:"docID"`
		EdinetCode           string      `json:"edinetCode"`
		SecCode              string      `json:"secCode"`
		Jcn                  string      `json:"JCN"`
		FilerName            string      `json:"filerName"`
		FundCode             interface{} `json:"fundCode"`
		OrdinanceCode        string      `json:"ordinanceCode"`
		FormCode             string      `json:"formCode"`
		DocTypeCode          string      `json:"docTypeCode"`
		PeriodStart          interface{} `json:"periodStart"`
		PeriodEnd            interface{} `json:"periodEnd"`
		SubmitDateTime       string      `json:"submitDateTime"`
		DocDescription       string      `json:"docDescription"`
		IssuerEdinetCode     interface{} `json:"issuerEdinetCode"`
		SubjectEdinetCode    interface{} `json:"subjectEdinetCode"`
		SubsidiaryEdinetCode interface{} `json:"subsidiaryEdinetCode"`
		CurrentReportReason  interface{} `json:"currentReportReason"`
		ParentDocID          interface{} `json:"parentDocID"`
		OpeDateTime          interface{} `json:"opeDateTime"`
		WithdrawalStatus     string      `json:"withdrawalStatus"`
		DocInfoEditStatus    string      `json:"docInfoEditStatus"`
		DisclosureStatus     string      `json:"disclosureStatus"`
		XbrlFlag             string      `json:"xbrlFlag"`
		PdfFlag              string      `json:"pdfFlag"`
		AttachDocFlag        string      `json:"attachDocFlag"`
		EnglishDocFlag       string      `json:"englishDocFlag"`
		CsvFlag              string      `json:"csvFlag"`
		LegalStatus          string      `json:"legalStatus"`
	} `json:"results"`
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
			p, err := parseParams()
			if err != nil {
				log.Fatal(err.Error())
			}
			fmt.Println("空のDBを作成します。")
			_, err = initDB()
			if err != nil {
				log.Fatal(err.Error())
			}

			parsedURL, err := url.Parse(EDINET_API_ENDPOINT)
			if err != nil {
				log.Fatal(err)
			}
			params := url.Values{}
			params.Add("Subscription-Key", viper.GetString("api.token"))
			params.Add("type", "2")
			for d := p.startDate; !d.After(p.endDate); d = d.AddDate(0, 0, 1) {
				params.Set("date", d.Format("2006-01-02"))
				parsedURL.RawQuery = params.Encode()
				resp, err := http.Get(parsedURL.String())
				if err != nil {
					log.Fatal(err)
				}
				defer resp.Body.Close()
				docs := Documents{}
				decoder := json.NewDecoder(resp.Body)
				if err := decoder.Decode(&docs); err != nil {
					log.Fatal(err)
				}
				fmt.Println("以下の日付のデータだよー")
				fmt.Println(d.String())
				fmt.Printf("%v", docs)
				fmt.Println()
				fmt.Println()
			}

		},
	}

	c.Flags().StringVarP(&strStartDate, "start-date", "s", "", "2023-12-12 の形式で追加してください")
	c.Flags().StringVarP(&strEndDate, "end-date", "e", "", "2023-12-12 の形式で追加してください")
	return c
}
