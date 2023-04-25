package main

import (
	"encoding/json"
	env "github.com/caitlinelfring/go-env-default"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
	"io"
	"net/http"
	"nginx_unit_exporter/connector"
)

type UnitStats struct {
	Connections  *ConnectionsState   `json:"connections,omitempty"`
	Requests     *InstanceRequests   `json:"requests,omitempty"`
	Applications map[string]*AppName `json:"applications,omitempty"`
}

type ConnectionsState struct {
	Accepted int `json:"accepted,omitempty"`
	Active   int `json:"active,omitempty"`
	Idle     int `json:"idle,omitempty"`
	Closed   int `json:"closed,omitempty"`
}
type InstanceRequests struct {
	Total int `json:"total,omitempty"`
}
type AppName struct {
	Processes map[string]int `json:"processes,omitempty"`
	Requests  *AppRequests   `json:"requests,omitempty"`
}

type AppRequests struct {
	Active int `json:"active,omitempty"`
}

// Descriptors used by the UnitCollector below.
var (
	unitInstanceRequestsTotalDesc = prometheus.NewDesc(
		"unit_instance_requests_total",
		"Total non-API requests during the instance’s lifetime.",
		[]string{"instance", "application"}, nil,
	)

	unitInstanceConnectionsAcceptedDesc = prometheus.NewDesc(
		"unit_instance_connections_accepted_total",
		"Total accepted connections during the instance’s lifetime.",
		[]string{"instance", "application"}, nil,
	)

	unitInstanceConnectionsActiveDesc = prometheus.NewDesc(
		"unit_instance_connections_active",
		"Current active connections for the instance",
		[]string{"instance", "application"}, nil,
	)

	unitInstanceConnectionsIdleDesc = prometheus.NewDesc(
		"unit_instance_connections_idle",
		"Current idle connections for the instance",
		[]string{"instance", "application"}, nil,
	)

	unitInstanceConnectionsClosedDesc = prometheus.NewDesc(
		"unit_instance_connections_closed_total",
		"Total closed connections during the instance’s lifetime",
		[]string{"instance", "application"}, nil,
	)

	unitApplicationRequestsActiveDesc = prometheus.NewDesc(
		"unit_application_requests_active",
		"Similar to /status/requests, but includes only the data for a specific app.",
		[]string{"instance", "application"}, nil,
	)

	unitApplicationProcessesDesc = prometheus.NewDesc(
		"unit_application_processes",
		"Current app processes.",
		[]string{"instance", "application", "state"}, nil,
	)
)

// UnitStatsCollector implements the Collector interface.
type UnitStatsCollector struct {
	UnitStats *UnitStats
}

func (sc UnitStatsCollector) Describe(ch chan<- *prometheus.Desc) {
	prometheus.DescribeByCollect(sc, ch)
}

func (sc UnitStatsCollector) Collect(ch chan<- prometheus.Metric) {
	resUnitMetrics, _ := sc.UnitStats.collectMetrics()

	unitInstanceRequestsTotal := resUnitMetrics.Requests.Total
	ch <- prometheus.MustNewConstMetric(
		unitInstanceRequestsTotalDesc,
		prometheus.CounterValue,
		float64(unitInstanceRequestsTotal),
		"unit", "",
	)

	unitInstanceConnectionsAccepted := resUnitMetrics.Connections.Accepted
	ch <- prometheus.MustNewConstMetric(
		unitInstanceConnectionsAcceptedDesc,
		prometheus.CounterValue,
		float64(unitInstanceConnectionsAccepted),
		"unit", "",
	)

	unitInstanceConnectionsActive := resUnitMetrics.Connections.Active
	ch <- prometheus.MustNewConstMetric(
		unitInstanceConnectionsActiveDesc,
		prometheus.GaugeValue,
		float64(unitInstanceConnectionsActive),
		"unit", "",
	)

	unitInstanceConnectionsIdle := resUnitMetrics.Connections.Idle
	ch <- prometheus.MustNewConstMetric(
		unitInstanceConnectionsIdleDesc,
		prometheus.GaugeValue,
		float64(unitInstanceConnectionsIdle),
		"unit", "",
	)

	unitInstanceConnectionsClosed := resUnitMetrics.Connections.Closed
	ch <- prometheus.MustNewConstMetric(
		unitInstanceConnectionsClosedDesc,
		prometheus.CounterValue,
		float64(unitInstanceConnectionsClosed),
		"unit", "",
	)

	unitApplicationRequestsActive := resUnitMetrics.Applications
	for application := range unitApplicationRequestsActive {
		ch <- prometheus.MustNewConstMetric(
			unitApplicationRequestsActiveDesc,
			prometheus.GaugeValue,
			float64(unitApplicationRequestsActive[application].Requests.Active),
			"unit", application,
		)
	}

	unitApplicationProcesses := resUnitMetrics.Applications
	for application := range unitApplicationProcesses {
		for state := range unitApplicationProcesses[application].Processes {
			ch <- prometheus.MustNewConstMetric(
				unitApplicationProcessesDesc,
				prometheus.GaugeValue,
				float64(unitApplicationProcesses[application].Processes[state]),
				"unit", application, state,
			)
		}
	}
}

// NewUnitStats first creates a Prometheus-ignorant UnitStats
// instance. Then, it creates a UnitStatsCollector for the just created
// UnitStats. Finally, it registers the UnitStatsCollector with a
// wrapping Registerer.
func NewUnitStats(reg prometheus.Registerer) *UnitStats {
	c := &UnitStats{}
	sc := UnitStatsCollector{UnitStats: c}
	prometheus.WrapRegistererWith(nil, reg).MustRegister(sc)
	return c
}

func main() {
	//#############################
	// Since we are dealing with custom Collector implementations, it might
	// be a good idea to try it out with a pedantic registry.
	reg := prometheus.NewPedanticRegistry()

	NewUnitStats(reg)

	// Add the standard process and Go metrics to the custom registry.
	reg.MustRegister(
		prometheus.NewProcessCollector(prometheus.ProcessCollectorOpts{}),
		// prometheus.NewGoCollector(),
	)

	http.Handle("/metrics", promhttp.HandlerFor(reg, promhttp.HandlerOpts{}))
	var metricsAddress = env.GetDefault("METRICS_LISTEN_ADDRESS", ":9095")
	log.Fatal(http.ListenAndServe(metricsAddress, nil))

}

func (stats *UnitStats) collectMetrics() (metrics *UnitStats, err error) {
	//############################
	// Here should be configured with flags or env
	// var network = "unix"
	// var address = "/tmp/unit-sock/control.unit.sock"
	var network = env.GetDefault("UNITD_CONTROL_NETWORK", "tcp")
	var address = env.GetDefault("UNITD_CONTROL_ADDRESS", ":8081")

	c := connector.NewConnection(network, address)
	res, err := c.Get("http://" + network + "/status")
	if err != nil {
		return metrics, err
	}
	body, err := io.ReadAll(res.Body)
	res.Body.Close()
	if res.StatusCode > 299 {
		log.Warnf("Response failed with status code: %d and\nbody: %s\n", res.StatusCode, body)
	}
	if err != nil {
		return metrics, err
	}

	err = json.Unmarshal(body, &metrics)
	if err != nil {
		return metrics, err
	}
	return metrics, nil

}
