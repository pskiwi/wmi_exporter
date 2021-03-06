// returns data points from Win32_PerfRawData_PerfOS_System class (undocumented)

package collector

import (
	"log"

	"github.com/StackExchange/wmi"
	"github.com/prometheus/client_golang/prometheus"
)

func init() {
	Factories["system"] = NewSystemCollector
}

// A PerfOSCollector is a Prometheus collector for WMI metrics
type SystemCollector struct {
	ContextSwitchesTotal     *prometheus.Desc
	ExceptionDispatchesTotal *prometheus.Desc
	ProcessorQueueLength     *prometheus.Desc
	SystemCallsTotal         *prometheus.Desc
	SystemUpTime             *prometheus.Desc
	Threads                  *prometheus.Desc
}

// NewSystemCollector ...
func NewSystemCollector() (Collector, error) {
	const subsystem = "system"

	return &SystemCollector{
		ContextSwitchesTotal: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, subsystem, "context_switches_total"),
			"PerfOS_System.ContextSwitchesPersec",
			nil,
			nil,
		),
		ExceptionDispatchesTotal: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, subsystem, "exception_dispatches_total"),
			"PerfOS_System.ExceptionDispatchesPersec",
			nil,
			nil,
		),
		ProcessorQueueLength: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, subsystem, "processor_queue_length"),
			"PerfOS_System.ProcessorQueueLength",
			nil,
			nil,
		),
		SystemCallsTotal: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, subsystem, "system_calls_total"),
			"PerfOS_System.SystemCallsPersec",
			nil,
			nil,
		),
		SystemUpTime: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, subsystem, "system_up_time"),
			"PerfOS_System.SystemUpTime",
			nil,
			nil,
		),
		Threads: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, subsystem, "threads"),
			"PerfOS_System.Threads",
			nil,
			nil,
		),
	}, nil
}

// Collect sends the metric values for each metric
// to the provided prometheus Metric channel.
func (c *SystemCollector) Collect(ch chan<- prometheus.Metric) error {
	if desc, err := c.collect(ch); err != nil {
		log.Println("[ERROR] failed collecting os metrics:", desc, err)
		return err
	}
	return nil
}

type Win32_PerfRawData_PerfOS_System struct {
	ContextSwitchesPersec     uint32
	ExceptionDispatchesPersec uint32
	ProcessorQueueLength      uint32
	SystemCallsPersec         uint32
	SystemUpTime              uint64
	Threads                   uint32
}

func (c *SystemCollector) collect(ch chan<- prometheus.Metric) (*prometheus.Desc, error) {
	var dst []Win32_PerfRawData_PerfOS_System
	if err := wmi.Query(wmi.CreateQuery(&dst, ""), &dst); err != nil {
		return nil, err
	}
	ch <- prometheus.MustNewConstMetric(
		c.ContextSwitchesTotal,
		prometheus.GaugeValue,
		float64(dst[0].ContextSwitchesPersec),
	)
	ch <- prometheus.MustNewConstMetric(
		c.ExceptionDispatchesTotal,
		prometheus.GaugeValue,
		float64(dst[0].ExceptionDispatchesPersec),
	)
	ch <- prometheus.MustNewConstMetric(
		c.ProcessorQueueLength,
		prometheus.GaugeValue,
		float64(dst[0].ProcessorQueueLength),
	)
	ch <- prometheus.MustNewConstMetric(
		c.SystemCallsTotal,
		prometheus.GaugeValue,
		float64(dst[0].SystemCallsPersec),
	)
	ch <- prometheus.MustNewConstMetric(
		c.SystemUpTime,
		prometheus.GaugeValue,
		float64(dst[0].SystemUpTime),
	)
	ch <- prometheus.MustNewConstMetric(
		c.Threads,
		prometheus.GaugeValue,
		float64(dst[0].Threads),
	)
	return nil, nil
}
