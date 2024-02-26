package main

import (
	"github.com/yuichi10/edinet_search/cmd"
)

//go:generate cp -r ./ui/edinet_search_ui/out ./cmd/api/out

func main() {
	cmd.Execute()
}
