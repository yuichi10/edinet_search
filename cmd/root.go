package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func newRootCmd() *cobra.Command {
	c := &cobra.Command{
		Use: "edinet",
		Short: "edinetから各会社の平均年収等を取ってきて表示します。",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("hello world")
		},
	}
	return c
}

func Execute() {
	rootCmd := newRootCmd()
	if err := rootCmd.Execute(); err != nil {
	  fmt.Fprintln(os.Stderr, err)
	  os.Exit(1)
	}
  }
