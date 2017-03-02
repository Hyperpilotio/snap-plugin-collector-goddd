package goddd

import (
	"fmt"
	"net/http"
	"time"

	dto "github.com/prometheus/client_model/go"
	"github.com/prometheus/common/expfmt"

	"github.com/intelsdi-x/snap/control/plugin"
	"github.com/intelsdi-x/snap/control/plugin/cpolicy"
	"github.com/intelsdi-x/snap/core"
	"github.com/intelsdi-x/snap/core/ctypes"
)

const (
	// Name of plugin
	Name = "goddd"
	// Version of plugin
	Version = 1
	// Type of plugin
	Type = plugin.CollectorPluginType
)

var (
	// make sure that we actually satisify requierd interface
	_ plugin.CollectorPlugin = (*Goddd)(nil)

	metricNames = []string{
		"avg",
		"min",
		"max",
		"median",
		"mdev",
		"loss",
	}
)

// Goddd struct
type Goddd struct {
}

// New return an instance of Goddd
func New() *Goddd {
	return &Goddd{}
}

func loadMetrics() map[string]*dto.MetricFamily {
	// FIXME url
	resp, err := http.Get("http://localhost:8080/metrics")
	if err != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()

	var parser expfmt.TextParser
	ans, err := parser.TextToMetricFamilies(resp.Body)

	if err != nil {
		fmt.Println("err is \n\n", err)
	}

	//ans, err := promql.ParseMetric()
	// for label, val := range ans {
	// 	fmt.Println(label)
	// 	fmt.Println(val.String(), "\n")
	// }

	return ans
}

// CollectMetrics collects metrics for testing
func (p *Goddd) CollectMetrics(mts []plugin.MetricType) ([]plugin.MetricType, error) {
	var err error

	conf := mts[0].Config().Table()
	hostname, ok := conf["hostname"]
	if !ok || hostname.(ctypes.ConfigValueStr).Value == "" {
		return nil, fmt.Errorf("hostname missing from config, %v", conf)
	}
	var timeout float64
	timeoutConf, ok := conf["timeout"]
	if !ok || timeoutConf.(ctypes.ConfigValueFloat).Value == 0 {
		timeout = 10.0
	} else {
		timeout = timeoutConf.(ctypes.ConfigValueFloat).Value
	}
	var count int
	countConf, ok := conf["count"]
	if !ok || countConf.(ctypes.ConfigValueInt).Value == 0 {
		count = 5
	} else {
		count = countConf.(ctypes.ConfigValueInt).Value
	}

	metrics, err := ping(hostname.(ctypes.ConfigValueStr).Value, count, timeout, mts)
	if err != nil {
		return nil, err
	}

	return metrics, nil
}

func ping(host string, count int, timeout float64, mts []plugin.MetricType) ([]plugin.MetricType, error) {
	check, err := NewGodddPingProbe(host, count, timeout)
	if err != nil {
		return nil, err
	}
	runTime := time.Now()
	result, err := check.Run()
	if err != nil {
		return nil, err
	}
	stats := make(map[string]float64)
	if result.Avg != nil {
		stats["avg"] = *result.Avg
	}
	if result.Min != nil {
		stats["min"] = *result.Min
	}
	if result.Max != nil {
		stats["max"] = *result.Max
	}
	if result.Median != nil {
		stats["median"] = *result.Median
	}
	if result.Mdev != nil {
		stats["mdev"] = *result.Mdev
	}
	if result.Loss != nil {
		stats["loss"] = *result.Loss
	}

	metrics := make([]plugin.MetricType, 0, len(stats))
	for _, m := range mts {
		stat := m.Namespace()[2].Value
		if value, ok := stats[stat]; ok {
			mt := plugin.MetricType{
				Data_:      value,
				Namespace_: core.NewNamespace("goddd", "ping", stat),
				Timestamp_: runTime,
				Version_:   m.Version(),
			}
			metrics = append(metrics, mt)
		}
	}

	return metrics, nil
}

//GetMetricTypes returns metric types for testing
func (p *Goddd) GetMetricTypes(cfg plugin.ConfigType) ([]plugin.MetricType, error) {
	mts := []plugin.MetricType{}
	for _, metricName := range metricNames {
		mts = append(mts, plugin.MetricType{
			Namespace_: core.NewNamespace("goddd", "ping", metricName),
		})
	}
	return mts, nil
}

//GetConfigPolicy returns a ConfigPolicyTree for testing
func (p *Goddd) GetConfigPolicy() (*cpolicy.ConfigPolicy, error) {
	c := cpolicy.New()
	rule0, _ := cpolicy.NewStringRule("hostname", true)
	rule1, _ := cpolicy.NewFloatRule("timeout", false, 10.0)
	rule2, _ := cpolicy.NewIntegerRule("count", false, 5)
	cp := cpolicy.NewPolicyNode()
	cp.Add(rule0)
	cp.Add(rule1)
	cp.Add(rule2)
	c.Add([]string{"goddd", "ping"}, cp)
	return c, nil
}

//Meta returns meta data for testing
func Meta() *plugin.PluginMeta {
	return plugin.NewPluginMeta(
		Name,
		Version,
		Type,
		[]string{plugin.SnapGOBContentType},
		[]string{plugin.SnapGOBContentType},
		plugin.Unsecure(true),
		plugin.ConcurrencyCount(5000),
	)
}
