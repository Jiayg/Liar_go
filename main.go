package main

import (
	_ "liar/internal/packed"

	"github.com/gogf/gf/v2/os/gctx"
	"liar/internal/cmd"
)

func main() {
	cmd.Main.Run(gctx.New())
}
