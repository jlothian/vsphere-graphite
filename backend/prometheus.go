package backend

// InitPrometheus : Set some channels to notify other theads when using Prometheus
import (
	"fmt"
	"log"

	"github.com/prometheus/client_golang/prometheus"
)

// Describe : Implementation of Prometheus Collector.Describe
func (backend *Config) Describe(ch chan<- *prometheus.Desc) {
	prometheus.NewGauge(prometheus.GaugeOpts{Name: "Dummy", Help: "Dummy"}).Describe(ch)
}

// Collect : Implementation of Prometheus Collector.Collect
func (backend *Config) Collect(ch chan<- prometheus.Metric) {

	log.Println("Requested Metrics!")

	*query <- true

	for {
		select {
		case point := <-*metrics:
			backend.PrometheusSend(ch, point)
		case <-*done:
			return
		}
	}
}

//PrometheusSend sends a point to prometheus
func (backend *Config) PrometheusSend(ch chan<- prometheus.Metric, point Point) {
	tags := point.GetTags(backend.NoArray, ",")
	labelNames := make([]string, len(tags))
	labelValues := make([]string, len(tags))
	i := 0
	for key, value := range tags {
		labelNames[i] = key
		labelValues[i] = value
		i++
	}
	key := fmt.Sprintf("%s_%s_%s_%s", backend.Prefix, point.Group, point.Counter, point.Rollup)
	desc := prometheus.NewDesc(key, "vSphere collected metric", labelNames, nil)
	metric, err := prometheus.NewConstMetric(desc, prometheus.GaugeValue, float64(point.Value), labelValues...)
	if err != nil {
		log.Println("Error creating prometheus metric")
	}
	ch <- metric
}
