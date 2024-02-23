package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/yuichi10/edinet_search/cmd/api"
	"github.com/yuichi10/edinet_search/cmd/createdb"
	"github.com/yuichi10/edinet_search/cmd/search"
)

func newRootCmd() *cobra.Command {
	cobra.OnInitialize(initConfig)

	c := &cobra.Command{
		Use:   "edinet",
		Short: "edinetから各会社の平均年収等を取ってきて表示します。",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("hello world")
		},
	}

	c.PersistentFlags().StringP("token", "t", "", "edinet api token. This cmd using v2")
	viper.BindPFlag("api.token", c.PersistentFlags().Lookup("token"))

	viper.SetDefault("EDINET_API_TOKEN", "")
	viper.BindEnv("api.token", "EDINET_API_TOKEN")

	c.AddCommand(createdb.NewCreateDBCmd())
	c.AddCommand(search.NewSearchCmd())
	c.AddCommand(api.New())

	return c
}

func initConfig() {
	viper.AutomaticEnv()
}

func Execute() {
	rootCmd := newRootCmd()
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
