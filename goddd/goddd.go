package goddd

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/golang/glog"
	"github.com/intelsdi-x/snap-plugin-lib-go/v1/plugin"
	dto "github.com/prometheus/client_model/go"
	"github.com/prometheus/common/expfmt"
)

var (
	vendor          = "hyperpilot"
	pluginName      = "goddd"
	pluginVersion   = 1
	nameSpacePrefix = []string{vendor, pluginName}
)

// GoCollector struct
type GoCollector struct {
}

// New return an instance of Goddd
func New() GoCollector {
	return GoCollector{}
}

// CollectMetrics will be called by Snap when a task that collects one of the metrics returned from this plugins
func (c GoCollector) CollectMetrics(mts []plugin.Metric) ([]plugin.Metric, error) {
	metrics := []plugin.Metric{}
	currentTime := time.Now()

	if len(mts) == 0 {
		return metrics, fmt.Errorf("array of metric type is empty\nPlease check GetMetricTypes()\n")
	}

	endpointConfig, err := mts[0].Config.GetString("endpoint")
	if err != nil {
		return metrics, fmt.Errorf("Unable to get endpoint config: " + err.Error())
	}

	endpoint := parseEndpoint(endpointConfig)

	metricFamilies, err := c.collect(endpoint)
	if err != nil {
		glog.Warningf("Unable to collect metrics, skipping to next cycle. endpoint: %s, error: %s", endpoint, err.Error())
		return metrics, nil
	}

	for idx, mt := range mts {
		mts[idx].Timestamp = currentTime
		ns := mt.Namespace.Strings()
		metricFamily := metricFamilies[ns[len(ns)-1]]

		metric := plugin.Metric{
			Namespace:   plugin.NewNamespace(ns...),
			Timestamp:   currentTime,
			Description: metricFamily.GetHelp(),
			Version:     int64(pluginVersion),
		}

		for _, metricOfGoddd := range metricFamily.GetMetric() {
			switch metricFamily.GetType() {

			case dto.MetricType_GAUGE:
				if strings.Contains(metricFamily.GetName(), "bytes") {
					metric.Unit = "B"
				}
				metric.Data = metricOfGoddd.GetGauge().GetValue()
				metric.Tags = getTagsOfMetric(metricOfGoddd)
				metrics = append(metrics, metric)

			case dto.MetricType_COUNTER:
				metric.Data = metricOfGoddd.GetCounter().GetValue()
				metric.Tags = getTagsOfMetric(metricOfGoddd)
				metrics = append(metrics, metric)

			case dto.MetricType_SUMMARY:
				summaryData, err := processSummaryMetric(metricOfGoddd)
				if err != nil {
					continue
				}
				for key, val := range summaryData {
					tags := getTagsOfMetric(metricOfGoddd)
					tags["summary"] = key
					metric.Tags = tags
					metric.Data = val
					metrics = append(metrics, metric)
				}

			}

		}
	}
	return metrics, nil
}

func getTagsOfMetric(metric *dto.Metric) map[string]string {
	tags := make(map[string]string)
	for _, label := range metric.GetLabel() {
		tags[label.GetName()] = label.GetValue()
	}
	return tags
}

func processSummaryMetric(metric *dto.Metric) (map[string]float64, error) {
	summary := make(map[string]float64)
	summary["count"] = float64(metric.GetSummary().GetSampleCount())
	summary["sum"] = float64(metric.GetSummary().GetSampleSum())
	summary["avg"] = summary["sum"] / summary["count"]

	for _, quantile := range metric.GetSummary().GetQuantile() {
		key := fmt.Sprintf("quantile_%d", int(quantile.GetQuantile()*100))
		summary[key] = quantile.GetValue()
	}

	return summary, nil
}

func downloadMetrics(url string) (io.Reader, error) {
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println(err)
		return nil, err
	} else if resp.StatusCode == http.StatusOK {
		defer resp.Body.Close()

		// Copy content from the body of http request
		buf := new(bytes.Buffer)
		buf.ReadFrom(resp.Body)
		b := buf.Bytes()
		httpBody := bytes.NewReader(b)

		return httpBody, nil
	} else {
		return nil, fmt.Errorf("Status code: %d Response: %v\n", resp.StatusCode, resp)
	}
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

func (c GoCollector) collect(endpoint string) (map[string]*dto.MetricFamily, error) {
	var httpBody io.Reader
	httpBody, err := downloadMetrics(endpoint)
	if err != nil {
		return nil, errors.New("Unable to download metrics: " + err.Error())
	}
	metricFamilies, err := parseMetrics(httpBody)
	if err != nil {
		return nil, errors.New("Unable to parse metrics: " + err.Error())
	}
	return metricFamilies, nil
}

//GetMetricTypes returns metric types for testing
func (c GoCollector) GetMetricTypes(cfg plugin.Config) ([]plugin.Metric, error) {
	mts := []plugin.Metric{}

	for _, val := range MetricList {
		mts = append(mts, plugin.Metric{
			// /hyperpilot/goddd/*
			Namespace: plugin.NewNamespace(nameSpacePrefix...).
				AddStaticElement(val),
			Version: int64(pluginVersion),
		})
	}

	return mts, nil
}

//GetConfigPolicy returns a ConfigPolicyTree for testing
func (c GoCollector) GetConfigPolicy() (plugin.ConfigPolicy, error) {
	policy := plugin.NewConfigPolicy()

	// name space
	configKey := nameSpacePrefix
	policy.AddNewStringRule(configKey,
		"endpoint",
		false,
		plugin.SetDefaultString("http://localhost:8080/metrics"))

	return *policy, nil
}

func parseEndpoint(address string) string {
	if strings.Contains(address, "/metrics") {
		return address
	}
	return address + "/metrics"
}
