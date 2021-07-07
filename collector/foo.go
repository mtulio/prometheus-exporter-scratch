package collector

import (
	"math/rand"

	"github.com/prometheus/client_golang/prometheus"
)

//Define a struct for you collector that contains pointers
//to prometheus descriptors for each metric you wish to expose.
//Note you can also include fields of other types if they provide utility
//but we just won't be exposing them as metrics.
type fooCollector struct {
	fooMetric *prometheus.Desc
}

const (
	fooCollectorSubsystem = "foo"
)

func init() {
	registerCollector("foo", defaultEnabled, NewFooCollector)
}

//You must create a constructor for you collector that
//initializes every descriptor and returns a pointer to the collector
func NewFooCollector() (Collector, error) {
	return &fooCollector{
		fooMetric: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, fooCollectorSubsystem, "metric_float"),
			"Shows a screatch metrics float value",
			[]string{"label1", "label2"}, nil,
		),
	}, nil
}

// Update implements Collector and exposes related metrics
func (c *fooCollector) Update(ch chan<- prometheus.Metric) error {
	if err := c.updateFoo(ch); err != nil {
		return err
	}
	return nil
}

func (c *fooCollector) updateFoo(ch chan<- prometheus.Metric) error {

	metricValue := (rand.Float64() * 5) + 5

	ch <- prometheus.MustNewConstMetric(
		c.fooMetric,
		prometheus.GaugeValue,
		metricValue,
		"value1", "value2",
	)

	return nil
}
