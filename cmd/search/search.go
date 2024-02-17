package search

import (
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3"
	"github.com/spf13/cobra"
	"github.com/yuichi10/edinet_search/db"
)

var companies []string

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

			fmt.Println(docs)
		},
	}

	c.Flags().StringSliceVarP(&companies, "companies", "c", make([]string, 0), "知りたい情報の会社の名前を書いていってください")

	return c
}
