package goddd

import (
	"io"
	"math"
	"os"
	"strings"
	"testing"

	"github.com/intelsdi-x/snap-plugin-lib-go/v1/plugin"

	. "github.com/smartystreets/goconvey/convey"
)

type MockMetricsDownloader struct {
	testData string
}

func (downloader *MockMetricsDownloader) SetTestData(dataName string) {
	downloader.testData = dataName
}

func (downloader *MockMetricsDownloader) GetMetricsReader(url string) (io.Reader, error) {
	if downloader.testData != "" {
		return strings.NewReader(TEST_DATA_SET[downloader.testData]), nil
	}
	return strings.NewReader(TEST_DATA_RAW), nil
}

func (downloader MockMetricsDownloader) GetEndpoint(config plugin.Config) (string, error) {
	return "test", nil
}
func deleteFileIfExist(filePath string) error {
	if _, err := os.Stat(filePath); err == nil {
		err = os.Remove(filePath)
		if err != nil {
			return err
		}
		return nil
	}
	return nil
}

func TestGodddPlugin(t *testing.T) {
	Convey("Create Goddd Collector", t, func() {
		collector := New()
		Convey("So Goddd collector should not be nil", func() {
			So(collector, ShouldNotBeNil)
		})
		Convey("collector.GetConfigPolicy() should return a config policy", func() {
			configPolicy, err := collector.GetConfigPolicy()
			Convey("So config policy should not be nil", func() {
				So(err, ShouldBeNil)
				So(configPolicy, ShouldNotBeNil)
				t.Log(configPolicy)
			})
			Convey("So config policy should be a policy.ConfigPolicy", func() {
				So(configPolicy, ShouldHaveSameTypeAs, plugin.ConfigPolicy{})
			})
		})
	})

	Convey("Test parsing metrics", t, func() {
		// collector := New()
		collector := &GodddCollector{
			Downloader: &MockMetricsDownloader{},
			cache:      NewCache(),
		}

		Convey("Goddd collect metrics should succesfully parse test metrics", func() {
			metricTypes, err := collector.GetMetricTypes(plugin.Config{})
			So(err, ShouldBeNil)
			metrics, err := collector.CollectMetrics(metricTypes)
			So(err, ShouldBeNil)

			Convey("Goddd collector should only store the diff of counter in a period of time", func() {
				Convey("Counter should be zero if metrics are the same ( no new traffic )", func() {
					deleteFileIfExist(cachePath)
					initCache()
					collector.cache = NewCache()
					mockDownloader := &MockMetricsDownloader{}
					mockDownloader.SetTestData("TEST_DATA_1")
					collector.Downloader = mockDownloader
					metrics, _ := collector.CollectMetrics(metricTypes)

					for _, metric := range metrics {
						ns := strings.Join(metric.Namespace.Strings(), "/")
						if ns == "hyperpilot/goddd/api_booking_service_request_count" {
							switch metric.Tags["method"] {
							case "list_cargos":
								So(metric.Data.(float64), ShouldBeZeroValue)
							case "list_locations":
								So(metric.Data.(float64), ShouldBeZeroValue)
							}
						}
					}
				})

				Convey("Counter should increase if application has new traffic ( increase 10,000 requests )", func() {
					mockDownloader := &MockMetricsDownloader{}
					mockDownloader.SetTestData("TEST_DATA_2")
					collector.Downloader = mockDownloader
					metrics, _ := collector.CollectMetrics(metricTypes)
					for _, metric := range metrics {
						ns := strings.Join(metric.Namespace.Strings(), "/")
						if ns == "hyperpilot/goddd/api_booking_service_request_count" {
							switch metric.Tags["method"] {
							case "list_cargos":
								So(metric.Data.(float64), ShouldEqual, 10000.0)
							case "list_locations":
								So(metric.Data.(float64), ShouldEqual, 10000.0)
							}
						}
					}
				})

				Convey("Total number of latency should increase if goddd has new traffic ( 20,000 requests and 40 seconds latency increased )", func() {
					deleteFileIfExist(cachePath)
					initCache()
					collector.cache = NewCache()
					mockDownloader := &MockMetricsDownloader{}
					mockDownloader.SetTestData("TEST_DATA_1")
					collector.Downloader = mockDownloader
					metrics, _ := collector.CollectMetrics(metricTypes)

					mockDownloader = &MockMetricsDownloader{}
					mockDownloader.SetTestData("TEST_DATA_2")
					collector.Downloader = mockDownloader
					metrics, _ = collector.CollectMetrics(metricTypes)

					for _, metric := range metrics {

						ns := strings.Join(metric.Namespace.Strings(), "/")
						summary, _ := metric.Tags["summary"]
						total, _ := metric.Tags["total"]

						if ns == "hyperpilot/goddd/api_booking_service_request_latency_microseconds" && total == "TOTAL" {
							switch summary {
							case "count":
								So(metric.Data.(float64), ShouldEqual, 20000)
							case "sum":
								So(metric.Data.(float64), ShouldEqual, 40)
							}
						}
					}
				})
			})

			Convey("Goddd collector should skip NaN metric values", func() {
				for _, metric := range metrics {
					So(math.IsNaN(metric.Data.(float64)), ShouldBeFalse)
				}
			})

			Convey("Goddd collector should calculate overall sum, count and average for metrics in MultiGroupMetricList", func() {
				deleteFileIfExist(cachePath)
				initCache()
				collector.cache = NewCache()
				mockDownloader := &MockMetricsDownloader{}
				mockDownloader.SetTestData("TEST_DATA_RAW")
				collector.Downloader = mockDownloader
				metrics, _ := collector.CollectMetrics(metricTypes)
				Convey("Overall sum, count, and average should be zero if goddd does not have new traffic ( no new requests )", func() {
					totalTagCounter := 0
					for _, metric := range metrics {
						ns := strings.Join(metric.Namespace.Strings(), "/")
						if total, _ := metric.Tags["total"]; total == "TOTAL" {
							switch ns {
							case "hyperpilot/goddd/api_booking_service_request_count":
								totalTagCounter++
								So(metric.Data.(float64), ShouldBeZeroValue)
							case "hyperpilot/goddd/api_booking_service_request_latency_microseconds":
								switch metric.Tags["summary"] {
								case "sum":
									totalTagCounter++
									So(metric.Data.(float64), ShouldBeZeroValue)
								case "count":
									totalTagCounter++
									So(metric.Data.(float64), ShouldBeZeroValue)
								case "avg":
									totalTagCounter++
									So(metric.Data.(float64), ShouldBeZeroValue)
								}
							}
						}
					}
					Convey("Goddd collector should collect exactly 4 overall summary metrics for the test metrics", func() {
						So(totalTagCounter, ShouldEqual, 4)
					})
				})
			})
		})
	})
}
