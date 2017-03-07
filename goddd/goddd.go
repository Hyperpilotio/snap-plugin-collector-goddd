package goddd

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"strings"
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
		metricOfGoddd := metricFamilies[ns[len(ns)-1]]

		metric := plugin.Metric{
			Namespace: plugin.NewNamespace(ns...),
			Timestamp: currentTime,
			Version:   int64(pluginVersion),
		}

		switch metricOfGoddd.GetType() {
		case dto.MetricType_GAUGE:
			if len(metricOfGoddd.GetMetric()) > 0 {
				if strings.Contains(metricOfGoddd.GetName(), "bytes") {
					metric.Unit = "B"
				}
				metric.Data = metricOfGoddd.GetMetric()[0].GetGauge().GetValue()
				metrics = append(metrics, metric)
			}
		case dto.MetricType_COUNTER:
			if len(metricOfGoddd.GetMetric()) > 0 {
				metric.Data = metricOfGoddd.GetMetric()[0].GetCounter().GetValue()
				metrics = append(metrics, metric)
			}
		case dto.MetricType_SUMMARY:
			if len(metricOfGoddd.GetMetric()) > 0 {

				if metric.Data, err = processSummaryMetric(metricOfGoddd.GetMetric()[0]); err != nil {
					metric.Data = ""
				}
				metrics = append(metrics, metric)
			}

		}
	}
	return metrics, nil
}

func processSummaryMetric(metric *dto.Metric) (string, error) {
	summaryStruct := Summary{}

	summaryStruct.SampleCount = metric.GetSummary().GetSampleCount()
	summaryStruct.SampleSum = metric.GetSummary().GetSampleSum()

	for _, quantile := range metric.GetSummary().GetQuantile() {
		switch quantile.GetQuantile() {
		case 0.5:
			summaryStruct.Quantile050 = quantile.GetQuantile()
		case 0.9:
			summaryStruct.Quantile090 = quantile.GetQuantile()
		case 0.99:
			summaryStruct.Quantile099 = quantile.GetQuantile()
		}
	}

	for _, label := range metric.GetLabel() {
		summaryStruct.Label = append(summaryStruct.Label, &LabelStruct{Name: label.GetName(), Value: label.GetValue()})
	}

	str, err := summaryStruct.MarshalJSON()
	if err != nil {
		return "", err
	}

	return string(str), nil
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

//GetMetricTypes returns metric types for testing
func (c GoCollector) GetMetricTypes(cfg plugin.Config) ([]plugin.Metric, error) {
	mts := []plugin.Metric{}
	metricFamilies, err := c.collect()
	if err != nil {
		return nil, fmt.Errorf("GoCollector.collect() called. Error: %s", err.Error())
	}

	//List of terminal metric names
	mList := make(map[string]bool)
	for _, val := range metricFamilies {
		// Keep it if not already seen before
		if !mList[val.GetName()] {
			mList[*val.Name] = true
			mts = append(mts, plugin.Metric{
				// /hyperpilot/goddd/*
				Namespace: plugin.NewNamespace(nameSpacePrefix...).
					AddStaticElement(val.GetName()),
				Description: val.GetHelp(),
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
