package main

import gop "github.com/lukebergen/gop/gopivot"

const version = "0.1"

func main() {
	gop.Init(version)
	gop.Exec()
}
