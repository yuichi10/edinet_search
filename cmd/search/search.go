package search

import (
	"fmt"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
	"github.com/yuichi10/edinet_search/db"
)

var companies []string
var verbose bool

func NewSearchCmd() *cobra.Command {
	c := &cobra.Command{
		Use:   "search",
		Short: "会社の情報を表示します。先にcreatedbコマンドを実施してから利用してください。",
		Run: func(cmd *cobra.Command, args []string) {
			db.OpenDB()
			docs, err := db.GetDocuments(companies)
			if err != nil {
				log.Fatal(err)
			}

			table := tablewriter.NewWriter(os.Stdout)
			table.SetHeader([]string{"会社名", "勤続年数", "平均年齢", "平均年収", "情報の追加日"})
			vTable := tablewriter.NewWriter(os.Stdout)
			vTable.SetHeader([]string{"会社名", "授業員情報"})
			vTable.SetAutoWrapText(false)

			for _, doc := range docs {
				data := []string{doc.FilerName, fmt.Sprintf("%s年",doc.AvgYearOfService), fmt.Sprintf("%s歳",doc.AvgAge), fmt.Sprintf("%s円",doc.AvgAnnualSalary), doc.SubmitDatetime}
				vData := []string{doc.FilerName, fmt.Sprintf("%s", doc.EmployeeInformation)}
				table.Append(data)
				vTable.Append(vData)
			}
			table.Render()
			if verbose {
				fmt.Println()
				fmt.Println("詳細情報")
				vTable.Render()
			}
		},
	}

	c.Flags().StringSliceVarP(&companies, "companies", "c", make([]string, 0), "知りたい情報の会社の名前を書いていってください")
	c.Flags().BoolVarP(&verbose, "verbose", "v", false, "授業員情報に関しても表示する。")

	return c
}
