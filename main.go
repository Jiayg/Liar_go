package main

import (
	_ "github.com/jiayg/liar/internal/packed"

	"github.com/gogf/gf/v2/os/gctx"
	"github.com/jiayg/liar/internal/cmd"
)

func main() {
	cmd.Main.Run(gctx.New())
}
