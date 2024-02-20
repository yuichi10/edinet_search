package search

import (
	"fmt"
	"log"
	"os"
	"strings"

	_ "github.com/mattn/go-sqlite3"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
	"github.com/yuichi10/edinet_search/db"
	"golang.org/x/term"
)

const EDINET_PDF_URL = "https://disclosure2dl.edinet-fsa.go.jp/searchdocument/pdf"

var companies []string
var salary string
var verbose bool

func makeNewLineText(text string, length int) string {
	replacer := strings.NewReplacer(
		"\t", "",
		"\r", "",
		"\n", "",
		" ", "",
	)
	text = replacer.Replace(text)

	runeText := []rune(text)

	var splits []string
	for start := 0; start < len(runeText); start += length {
		end := start + length
		if end > len(runeText) {
			end = len(runeText)
		}
		split := runeText[start:end]
		splits = append(splits, string(split))
	}

	return strings.Join(splits, "\n")
}

func NewSearchCmd() *cobra.Command {
	c := &cobra.Command{
		Use:   "search",
		Short: "会社の情報を表示します。先にcreatedbコマンドを実施してから利用してください。",
		Run: func(cmd *cobra.Command, args []string) {
			db.OpenDB()
			docs, err := db.GetCompanies(companies, salary)
			if err != nil {
				log.Fatal(err)
			}

			width, _, err := term.GetSize(int(os.Stdout.Fd()))
			width = width / 3 // マルチバイトの日本語を考えて /2　あとは会社名分を考えて合わせて/3くらいをしている。
			if err != nil {
				width = 30
			}

			table := tablewriter.NewWriter(os.Stdout)
			// table.SetHeader([]string{"会社名", "勤続年数", "平均年齢", "平均年収", "従業員数", "情報の追加日"})
			// table.SetAutoFormatHeaders(false)
			// table.EnableBorder(false)
			table.SetBorders(tablewriter.Border{Left: true, Top: false, Right: true, Bottom: false})
			vTable := tablewriter.NewWriter(os.Stdout)
			vTable.SetHeader([]string{"会社名", "授業員情報"})
			vTable.SetRowLine(true)
			vTable.SetRowSeparator("-")

			table.Append([]string{fmt.Sprintf("\x1b]8;;%s\x1b\\%s\x1b]8;;\x1b\\", fmt.Sprintf("%s/%s.pdf",EDINET_PDF_URL, "DUMMY012"), "会社名"), "勤続年数", "平均年齢", "平均年収", "従業員数", "情報の追加日"})
			for _, doc := range docs {
				pdfURL := fmt.Sprintf("%s/%s.pdf", EDINET_PDF_URL, doc.DocID)
				data := []string{fmt.Sprintf("\x1b]8;;%s\x1b\\%s\x1b]8;;\x1b\\", pdfURL, doc.FilerName), fmt.Sprintf("%s年", doc.AvgYearOfService), fmt.Sprintf("%s歳", doc.AvgAge), fmt.Sprintf("%s円", doc.AvgAnnualSalary), fmt.Sprintf("%s人", doc.NumberOfEmployees), doc.SubmitDatetime}
				vData := []string{doc.FilerName, makeNewLineText(doc.EmployeeInformation, width)}
				table.Append(data)
				vTable.Append(vData)
			}
			table.Render()
			if verbose {
				fmt.Println()
				fmt.Println("詳細情報")
				vTable.Render()
			}

			// fmt.Println(makeNewLineText(docs[0].EmployeeInformation, width))
		},
	}

	c.Flags().StringSliceVarP(&companies, "companies", "c", make([]string, 0), "知りたい情報の会社の名前を書いていってください")
	c.Flags().StringVarP(&salary, "salary", "s", "", "平均年収がいくら以上の会社を検索したいか記入してください")
	c.Flags().BoolVarP(&verbose, "verbose", "v", false, "授業員情報に関しても表示する。")

	return c
}
