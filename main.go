package main

import (
	"github.com/hyperpilotio/snap-plugin-collector-goddd/goddd"
	"github.com/intelsdi-x/snap-plugin-lib-go/v1/plugin"
)

const (
	pluginName    = "snap-plugin-collector-goddd"
	pluginVersion = 1
)

func main() {
	plugin.StartCollector(goddd.New(), pluginName, pluginVersion)
}
