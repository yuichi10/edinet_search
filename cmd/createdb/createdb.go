package createdb

import (
	"archive/zip"
	"bufio"
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"path"
	"strings"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/yuichi10/edinet_search/db"
	"golang.org/x/text/encoding/unicode"
	"golang.org/x/text/transform"
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


documentsテーブル
実際にedinetから取得してきたドキュメントの情報を保管。ドキュメントと行っているが、有価証券の情報だけを保管する予定.
docIDをprimayにしておけば良さげ。
ordinanceCodeが010、form_codeが030000のデータのみをまずは取るようにしてみる。
- docID
- parentDocID
- filerName
- submitDateTime
- docDescription

実際にファイルから取得したデータをDBに保存するような機能がほしい。
というか最初にすべてのmetadataを取ってきたときについでにすべてのデータをパースしてDBにデータを追加すれば良さそう。
とりあえず最初はすべてのデータを保存するようにしてしまう。
securities テーブル
doc_id 上のテーブルとの紐づけ(実際一つのテーブルに入れてしまってもいいのかもしれない。とりあえずのものだし、別テーブルで作ってしまう。)
number_of_employees 従業員数（人）
avg_age 平均年齢（歳）
avg_years_of_service 平均勤続年数（年）
avg_annual_salary 平均年間給与（円）
employee_information_text InformationAboutEmployeesTextBlockの情報をそのままいれる。

TODO
どこかで売上高と、経営利益を取得して、利益率も表示させたい。
*/

const EDINET_DOCUMENT_META_API_ENDPOINT = "https://api.edinet-fsa.go.jp/api/v2/documents.json"
const EDINET_DOCUMENT_API_ENDPOINT = "https://api.edinet-fsa.go.jp/api/v2/documents/"

// args
var strStartDate, strEndDate string

type inputParams struct {
	startDate time.Time
	endDate   time.Time
}

type DocumentMeta struct {
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
}

type DocumentsMeta struct {
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
	Results []DocumentMeta `json:"results"`
}

type SecuritiesInfo struct {
	NumberOfEmployees   string
	AvgAge              string
	AvgYearOfService    string
	AvgAnnualSalary     string
	EmployeeInformation string
}

func initDB() {
	err := db.DeleteDB()
	if err != nil {
		log.Fatal(err)
	}

	err = db.OpenDB()
	if err != nil {
		log.Fatal(err)
	}

	err = db.CreateCompaniesTable()
	if err != nil {
		log.Fatal(err)
	}
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

func getDocumentMetaInfo(d time.Time) (*DocumentsMeta, error) {
	parsedURL, err := url.Parse(EDINET_DOCUMENT_META_API_ENDPOINT)
	if err != nil {
		log.Fatal(err)
	}
	params := url.Values{}
	params.Add("Subscription-Key", viper.GetString("api.token"))
	params.Add("type", "2")
	params.Add("date", d.Format("2006-01-02"))
	parsedURL.RawQuery = params.Encode()
	resp, err := http.Get(parsedURL.String())
	if err != nil {
		return nil, fmt.Errorf("ドキュメントデータの取得に失敗しました %s", err)

	}
	defer resp.Body.Close()
	docs := DocumentsMeta{}
	decoder := json.NewDecoder(resp.Body)
	if err := decoder.Decode(&docs); err != nil {
		return nil, fmt.Errorf("取得したドキュメントデータのパースに失敗しました %s", err)
	}
	return &docs, nil
}

func getSecuritiesInfo(meta *DocumentMeta) (*SecuritiesInfo, error) {
	docParsedURL, err := url.Parse(EDINET_DOCUMENT_API_ENDPOINT)
	if err != nil {
		// return nil, fmt.Errorf("ドキュメントURLのパースに失敗しました %s", err)
		log.Fatalln(err)
	}
	docParsedURL.Path = path.Join(docParsedURL.Path, meta.DocID)
	params := url.Values{}
	params.Add("Subscription-Key", viper.GetString("api.token"))
	params.Add("type", "5")
	docParsedURL.RawQuery = params.Encode()
	resp, err := http.Get(docParsedURL.String())
	if err != nil {
		return nil, fmt.Errorf("ドキュメントの取得に失敗しました %s", err)
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("docID: %s, 会社名: %s のドキュメントファイルの取得に失敗しました", meta.DocID, meta.FilerName)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("ドキュメントファイルの読み込みに失敗しました %s", err)
	}

	zipReader, err := zip.NewReader(bytes.NewReader(body), int64(len(body)))
	if err != nil {
		return nil, fmt.Errorf("docID: %s, 会社名: %s のドキュメントファイルの解答に失敗しました %s", meta.DocID, meta.FilerName, err)
	}

	securInfo := &SecuritiesInfo{}
	for _, file := range zipReader.File {
		if strings.HasPrefix(file.Name, "XBRL_TO_CSV/jpcrp") {
			f, err := file.Open()
			if err != nil {
				return nil, fmt.Errorf("zip内のファイルの読み込みに失敗しました。 %s", err)
			}
			defer f.Close()
			decoder := unicode.UTF16(unicode.LittleEndian, unicode.UseBOM).NewDecoder()
			reader := csv.NewReader(transform.NewReader(bufio.NewReader(f), decoder))
			reader.Comma = '\t'
			for {
				record, err := reader.Read()
				if err == io.EOF {
					break
				}
				if err != nil {
					fmt.Printf("CSVファイルを読むのに失敗しました。 %s\n", err)
					break
				}

				switch record[0] {
				case "jpcrp_cor:AverageLengthOfServiceYearsInformationAboutReportingCompanyInformationAboutEmployees":
					// 平均継続年数
					securInfo.AvgYearOfService = record[len(record)-1]
				case "jpcrp_cor:AverageAgeYearsInformationAboutReportingCompanyInformationAboutEmployees":
					// 平均年齢
					securInfo.AvgAge = record[len(record)-1]
				case "jpcrp_cor:AverageAnnualSalaryInformationAboutReportingCompanyInformationAboutEmployees":
					// 平均年間給与
					securInfo.AvgAnnualSalary = record[len(record)-1]
				case "jpcrp_cor:InformationAboutEmployeesTextBlock":
					securInfo.EmployeeInformation = record[len(record)-1]
				case "jpcrp_cor:NumberOfEmployees":
					// 従業員数 (当社)
					if record[2] == "CurrentYearInstant_NonConsolidatedMember" {
						securInfo.NumberOfEmployees = record[len(record)-1]
					}
				}
			}
		}
	}

	return securInfo, nil
}

func convertDataToDBCompanies(meta *DocumentMeta, secInfo *SecuritiesInfo) db.Companies {
	return db.Companies{
		DocID:               meta.DocID,
		SecCode:             meta.SecCode,
		FilerName:           meta.FilerName,
		DocDescription:      meta.DocDescription,
		SubmitDatetime:      meta.SubmitDateTime,
		NumberOfEmployees:   secInfo.NumberOfEmployees,
		AvgAge:              secInfo.AvgAge,
		AvgYearOfService:    secInfo.AvgYearOfService,
		AvgAnnualSalary:     secInfo.AvgAnnualSalary,
		EmployeeInformation: secInfo.EmployeeInformation,
	}
}

func insertSecuritiesData(p inputParams) {
	for d := p.startDate; !d.After(p.endDate); d = d.AddDate(0, 0, 1) {
		meta, err := getDocumentMetaInfo(d)
		if err != nil {
			fmt.Printf("%s のデータの取得に失敗しました。 %s", d, err)
			continue
		}

		for _, m := range meta.Results {
			if m.OrdinanceCode == "010" && m.FormCode == "030000" {
				fmt.Print(".")
				docInfo, err := getSecuritiesInfo(&m)
				if err != nil {
					fmt.Printf("docID: %s, 会社名: %s 有価証券情報の取得に失敗しました。 %s", m.DocID, m.FilerName, err)
					continue
				}
				dbData := convertDataToDBCompanies(&m, docInfo)
				err = db.InsertCompanies(dbData)
				if err != nil {
					fmt.Printf("docID: %s, 会社名: %s のデータ保存に失敗しました。 %s", m.DocID, m.FilerName, err)
					continue
				}
			}
		}
	}
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
			initDB()

			fmt.Println("ドキュメントの情報をDBに書き込んでいます.")
			insertSecuritiesData(p)
		},
	}

	c.Flags().StringVarP(&strStartDate, "start-date", "s", "", "2023-12-12 の形式で追加してください")
	c.Flags().StringVarP(&strEndDate, "end-date", "e", "", "2023-12-12 の形式で追加してください")
	return c
}
