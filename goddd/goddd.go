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

	endpoint, err := parseEndpoint(mts[0].Config.GetString("endpoint"))
	if err != nil {
		return metrics, fmt.Errorf("Error on mts[0].Config.GetString(endpoint)\nError: %s\nendpoint: %s", err.Error(), endpoint)
	}

	metricFamilies, err := c.collect(endpoint)
	if err != nil {
		return metrics, fmt.Errorf("Error on GoCollector.collect():\n%s\nendpoint: %s\n", err.Error(), endpoint)

	}

	for idx, mt := range mts {
		mts[idx].Timestamp = currentTime
		ns := mt.Namespace.Strings()
		metricOfGoddd := metricFamilies[ns[len(ns)-1]]

		metric := plugin.Metric{
			Namespace:   plugin.NewNamespace(ns...),
			Timestamp:   currentTime,
			Description: metricOfGoddd.GetHelp(),
			Version:     int64(pluginVersion),
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
	} else if resp.StatusCode == http.StatusOK {
		defer resp.Body.Close()

		// Copy content from the body of http request
		buf := new(bytes.Buffer)
		buf.ReadFrom(resp.Body)
		b := buf.Bytes()
		httpBody := bytes.NewReader(b)

		return httpBody, nil
	} else {
		return nil, fmt.Errorf("Status code:%d Response: %v\n", resp.StatusCode, resp)
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
	metricFamilies, err := parseMetrics(httpBody)
	if err != nil {
		return nil, err
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

func parseEndpoint(address string, err error) (string, error) {
	if err != nil {
		return address, err
	}

	if strings.Contains(address, "/metrics") {
		return address, err
	}
	return address + "/metrics", err
}
