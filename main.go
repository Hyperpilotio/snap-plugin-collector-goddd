package main

import (
	// "github.com/intelsdi-x/snap/control/plugin"
	"github.com/intelsdi-x/snap-plugin-lib-go/v1/plugin"
	"github.com/swhsiang/snap-plugin-collector-goddd/goddd"
)

const (
	pluginName    = "snap-plugin-collector-goddd"
	pluginVersion = 1
)

func main() {
	plugin.StartCollector(goddd.New(), pluginName, pluginVersion)
}
