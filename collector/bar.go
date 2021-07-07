package collector

import (
	"math/rand"

	"github.com/prometheus/client_golang/prometheus"
)

type barCollector struct {
	barMetric *prometheus.Desc
}

const (
	barCollectorSubsystem = "bar"
)

func init() {
	registerCollector("bar", defaultEnabled, NewBarCollector)
}

func NewBarCollector() (Collector, error) {
	return &barCollector{
		barMetric: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, barCollectorSubsystem, "metric_float"),
			"Shows a scretch metrics int value",
			[]string{"label1", "label2"}, nil,
		),
	}, nil
}

func (c *barCollector) Update(ch chan<- prometheus.Metric) error {
	if err := c.updateBar(ch); err != nil {
		return err
	}
	return nil
}

func (c *barCollector) updateBar(ch chan<- prometheus.Metric) error {

	metricValue := (rand.Float64() * 5) + 5

	ch <- prometheus.MustNewConstMetric(
		c.barMetric,
		prometheus.GaugeValue,
		metricValue,
		"value1", "value2",
	)

	return nil
}
