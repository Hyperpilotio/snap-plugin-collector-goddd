package main

import (
	"os"

	"github.com/intelsdi-x/snap/control/plugin"
	"github.com/swhsiang/snap-plugin-collector-goddd/goddd"
)

func main() {

	plugin.Start(
		goddd.Meta(),
		goddd.New(),
		os.Args[1],
	)
}
