package goddd

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"time"

	dto "github.com/prometheus/client_model/go"
	"github.com/prometheus/common/expfmt"

	"github.com/intelsdi-x/snap-plugin-lib-go/v1/plugin"
	//"github.com/intelsdi-x/snap/control/plugin"
	//"github.com/intelsdi-x/snap/control/plugin/cpolicy"
	// "github.com/intelsdi-x/snap/core"
	//"github.com/intelsdi-x/snap/core/ctypes"
)

var (
	// make sure that we actually satisify requierd interface
	//_ plugin.GoCollectorPlugin = (*Collector)(nil)

	vendor          = "hyperpilot"
	pluginName      = "goddd"
	pluginVersion   = 1
	nameSpacePrefix = []string{vendor, pluginName}
)

// GoCollector struct
type GoCollector struct {
	URL string
}

// New return an instance of Goddd
func New(url string) GoCollector {
	return GoCollector{URL: url}
}

// CollectMetrics will be called by Snap when a task that collects one of the metrics returned from this plugins
func (c GoCollector) CollectMetrics(mts []plugin.Metric) ([]plugin.Metric, error) {
	metrics := []plugin.Metric{}
	currentTime := time.Now()

	metricFamilies, err := c.collect()
	if err != nil {
		return metrics, fmt.Errorf("Error on GoCollector.collect():\n%s", err.Error())

	}

	for idx, mt := range mts {
		mts[idx].Timestamp = currentTime
		ns := mt.Namespace.Strings()

		switch ns[len(ns)-1] {
		case "go_memstats_mcache_sys_bytes":
			{
				metric := plugin.Metric{
					Namespace: plugin.NewNamespace(ns...),
					Data:      metricFamilies[ns[len(ns)-1]].GetMetric()[0].GetGauge().GetValue(),
					Timestamp: currentTime,
					Version:   int64(pluginVersion),
					Unit:      "B",
				}
				fmt.Printf("[YO] go_memstats_mcache_sys_bytes got\n%v", metric)
				metrics = append(metrics, metric)
			}
		default:
			{
				return nil, fmt.Errorf("Invalid metric: %v", ns)
			}
		}
	}
	return metrics, nil
}

func downloadMetrics(url string) (io.Reader, error) {
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	defer resp.Body.Close()

	// Copy content from the body of http request
	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	b := buf.Bytes()
	httpBody := bytes.NewReader(b)

	return httpBody, nil
}

func parseMetrics(httpBody io.Reader) (map[string]*dto.MetricFamily, error) {
	var parser expfmt.TextParser
	metricFamilies, err := parser.TextToMetricFamilies(httpBody)
	if err != nil {
		fmt.Println(err)
		return make(map[string]*dto.MetricFamily), err
	}
	return metricFamilies, nil
}

func (c GoCollector) collect() (map[string]*dto.MetricFamily, error) {
	var httpBody io.Reader
	httpBody, err := downloadMetrics(c.URL)
	metricFamilies, err := parseMetrics(httpBody)
	if err != nil {
		return nil, err
	}
	return metricFamilies, nil
}

// func ping(host string, count int, timeout float64, mts []plugin.Metric) ([]plugin.Metric, error) {
// 	check, err := NewGodddPingProbe(host, count, timeout)
// 	if err != nil {
// 		return nil, err
// 	}
// 	runTime := time.Now()
// 	result, err := check.Run()
// 	if err != nil {
// 		return nil, err
// 	}
// 	stats := make(map[string]float64)
// 	if result.Avg != nil {
// 		stats["avg"] = *result.Avg
// 	}
// 	if result.Min != nil {
// 		stats["min"] = *result.Min
// 	}
// 	if result.Max != nil {
// 		stats["max"] = *result.Max
// 	}
// 	if result.Median != nil {
// 		stats["median"] = *result.Median
// 	}
// 	if result.Mdev != nil {
// 		stats["mdev"] = *result.Mdev
// 	}
// 	if result.Loss != nil {
// 		stats["loss"] = *result.Loss
// 	}

// 	metrics := make([]plugin.Metric, 0, len(stats))
// 	for _, m := range mts {
// 		stat := m.Namespace()[2].Value
// 		if value, ok := stats[stat]; ok {
// 			mt := plugin.MetricType{
// 				Data_:      value,
// 				Namespace_: core.NewNamespace("goddd", "ping", stat),
// 				Timestamp_: runTime,
// 				Version_:   m.Version(),
// 			}
// 			metrics = append(metrics, mt)
// 		}
// 	}

// 	return metrics, nil
// }

func parseMetricKey(val dto.MetricFamily) {

}

//GetMetricTypes returns metric types for testing
func (c GoCollector) GetMetricTypes(cfg plugin.Config) ([]plugin.Metric, error) {
	mts := []plugin.Metric{}
	metricFamilies, err := c.collect()
	if err != nil {
		return nil, err
	}

	//List of terminal metric names
	mList := make(map[string]bool)
	for _, val := range metricFamilies {
		// Keep it if not already seen before
		if !mList[*val.Name] {
			mList[*val.Name] = true
			mts = append(mts, plugin.Metric{
				// /hyperpilot/goddd/*
				Namespace: plugin.NewNamespace(nameSpacePrefix...).
					AddStaticElement(*val.Name),
				Description: *val.Help,
				Tags:        map[string]string{"type": val.GetType().String()},
				Version:     int64(pluginVersion),
			})
		}
	}
	return mts, nil
}

//GetConfigPolicy returns a ConfigPolicyTree for testing
func (c GoCollector) GetConfigPolicy() (plugin.ConfigPolicy, error) {
	policy := plugin.NewConfigPolicy()
	return *policy, nil
}
