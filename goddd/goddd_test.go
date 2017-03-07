package goddd

import (
	"testing"

	"github.com/intelsdi-x/snap-plugin-lib-go/v1/plugin"

	. "github.com/smartystreets/goconvey/convey"
)

func TestPingPlugin(t *testing.T) {
	Convey("Create GoCollector Collector", t, func() {
		collector := New()
		Convey("So Ping collector should not be nil", func() {
			So(collector, ShouldNotBeNil)
		})
		Convey("So ping collector should be of Ping type", func() {
			So(collector, ShouldHaveSameTypeAs, GoCollector{})
		})
		Convey("collector.GetConfigPolicy() should return a config policy", func() {
			configPolicy, _ := collector.GetConfigPolicy()
			Convey("So config policy should not be nil", func() {
				So(configPolicy, ShouldNotBeNil)
				t.Log(configPolicy)
			})
			Convey("So config policy should be a policy.ConfigPolicy", func() {
				So(configPolicy, ShouldHaveSameTypeAs, plugin.ConfigPolicy{})
			})
		})
	})
}
