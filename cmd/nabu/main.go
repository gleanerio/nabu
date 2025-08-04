package main

import (
	"github.com/gleanerio/nabu/internal/common"
	"github.com/gleanerio/nabu/pkg/cli"
)

func init() {
	common.InitLogging()
}

func main() {
	cli.Execute()
}
