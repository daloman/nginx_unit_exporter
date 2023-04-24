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

/*
package main

import (
	"log"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// ClusterManager is an example for a system that might have been built without
// Prometheus in mind. It models a central manager of jobs running in a
// cluster. Thus, we implement a custom Collector called
// ClusterManagerCollector, which collects information from a ClusterManager
// using its provided methods and turns them into Prometheus Metrics for
// collection.
//
// An additional challenge is that multiple instances of the ClusterManager are
// run within the same binary, each in charge of a different zone. We need to
// make use of wrapping Registerers to be able to register each
// ClusterManagerCollector instance with Prometheus.
type ClusterManager struct {
	Zone string
	// Contains many more fields not listed in this example.
}

// ReallyExpensiveAssessmentOfTheSystemState is a mock for the data gathering a
// real cluster manager would have to do. Since it may actually be really
// expensive, it must only be called once per collection. This implementation,
// obviously, only returns some made-up data.
func (c *ClusterManager) ReallyExpensiveAssessmentOfTheSystemState() (
	oomCountByHost map[string]int, ramUsageByHost map[string]float64,
) {
	// Just example fake data.
	oomCountByHost = map[string]int{
		"foo.example.org": 42,
		"bar.example.org": 2001,
	}
	ramUsageByHost = map[string]float64{
		"foo.example.org": 6.023e23,
		"bar.example.org": 3.14,
	}
	return
}

// ClusterManagerCollector implements the Collector interface.
type ClusterManagerCollector struct {
	ClusterManager *ClusterManager
}

// Descriptors used by the ClusterManagerCollector below.
var (
	oomCountDesc = prometheus.NewDesc(
		"clustermanager_oom_crashes_total",
		"Number of OOM crashes.",
		[]string{"host"}, nil,
	)
	ramUsageDesc = prometheus.NewDesc(
		"clustermanager_ram_usage_bytes",
		"RAM usage as reported to the cluster manager.",
		[]string{"host"}, nil,
	)
)

// Describe is implemented with DescribeByCollect. That's possible because the
// Collect method will always return the same two metrics with the same two
// descriptors.
func (cc ClusterManagerCollector) Describe(ch chan<- *prometheus.Desc) {
	prometheus.DescribeByCollect(cc, ch)
}

// Collect first triggers the ReallyExpensiveAssessmentOfTheSystemState. Then it
// creates constant metrics for each host on the fly based on the returned data.
//
// Note that Collect could be called concurrently, so we depend on
// ReallyExpensiveAssessmentOfTheSystemState to be concurrency-safe.
func (cc ClusterManagerCollector) Collect(ch chan<- prometheus.Metric) {
	oomCountByHost, ramUsageByHost := cc.ClusterManager.ReallyExpensiveAssessmentOfTheSystemState()
	for host, oomCount := range oomCountByHost {
		ch <- prometheus.MustNewConstMetric(
			oomCountDesc,
			prometheus.CounterValue,
			float64(oomCount),
			host,
		)
	}
	for host, ramUsage := range ramUsageByHost {
		ch <- prometheus.MustNewConstMetric(
			ramUsageDesc,
			prometheus.GaugeValue,
			ramUsage,
			host,
		)
	}
}

// NewClusterManager first creates a Prometheus-ignorant ClusterManager
// instance. Then, it creates a ClusterManagerCollector for the just created
// ClusterManager. Finally, it registers the ClusterManagerCollector with a
// wrapping Registerer that adds the zone as a label. In this way, the metrics
// collected by different ClusterManagerCollectors do not collide.
func NewClusterManager(zone string, reg prometheus.Registerer) *ClusterManager {
	c := &ClusterManager{
		Zone: zone,
	}
	cc := ClusterManagerCollector{ClusterManager: c}
	prometheus.WrapRegistererWith(prometheus.Labels{"zone": zone}, reg).MustRegister(cc)
	return c
}

func main() {
	// Since we are dealing with custom Collector implementations, it might
	// be a good idea to try it out with a pedantic registry.
	reg := prometheus.NewPedanticRegistry()

	// Construct cluster managers. In real code, we would assign them to
	// variables to then do something with them.
	NewClusterManager("db", reg)
	NewClusterManager("ca", reg)

	// Add the standard process and Go metrics to the custom registry.
	reg.MustRegister(
		prometheus.NewProcessCollector(prometheus.ProcessCollectorOpts{}),
		prometheus.NewGoCollector(),
	)

	http.Handle("/metrics", promhttp.HandlerFor(reg, promhttp.HandlerOpts{}))
	log.Fatal(http.ListenAndServe(":8080", nil))
}
*/

type UnitStats struct {
	Connections  map[string]int                       `json:"connections,omitempty"`
	Requests     map[string]int                       `json:"requests,omitempty"`
	Applications map[string]map[string]map[string]int `json:"applications,omitempty"`
}

// Descriptors used by the UnitCollector below.
var (
	unitInstanceRequestsTotalDesc = prometheus.NewDesc(
		"unit_instance_requests_total",
		"Total non-API requests during the instance’s lifetime.",
		[]string{"instance", "application"}, nil,
	)
	unitInstanceConnectionsAcceptedDesc = prometheus.NewDesc(
		"unit_instance_connections_accepted",
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
		"unit_instance_connections_closed",
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

	unitInstanceRequestsTotal := resUnitMetrics.Requests["total"]
	ch <- prometheus.MustNewConstMetric(
		unitInstanceRequestsTotalDesc,
		prometheus.CounterValue,
		float64(unitInstanceRequestsTotal),
		"unit", "",
	)

	unitInstanceConnectionsAccepted := resUnitMetrics.Connections["accepted"]
	ch <- prometheus.MustNewConstMetric(
		unitInstanceConnectionsAcceptedDesc,
		prometheus.CounterValue,
		float64(unitInstanceConnectionsAccepted),
		"unit", "",
	)

	unitInstanceConnectionsActive := resUnitMetrics.Connections["active"]
	ch <- prometheus.MustNewConstMetric(
		unitInstanceConnectionsActiveDesc,
		prometheus.GaugeValue,
		float64(unitInstanceConnectionsActive),
		"unit", "",
	)

	unitInstanceConnectionsIdle := resUnitMetrics.Connections["idle"]
	ch <- prometheus.MustNewConstMetric(
		unitInstanceConnectionsIdleDesc,
		prometheus.GaugeValue,
		float64(unitInstanceConnectionsIdle),
		"unit", "",
	)

	unitInstanceConnectionsClosed := resUnitMetrics.Connections["closed"]
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
			float64(unitApplicationRequestsActive[application]["requests"]["active"]),
			"unit", application,
		)
	}

	unitApplicationProcesses := resUnitMetrics.Applications
	for application := range unitApplicationProcesses {
		for state := range unitApplicationProcesses[application]["processes"] {
			ch <- prometheus.MustNewConstMetric(
				unitApplicationProcessesDesc,
				prometheus.GaugeValue,
				float64(unitApplicationProcesses[application]["processes"][state]),
				"unit", application, state,
			)
		}
	}
}

// NewUnitStats first creates a Prometheus-ignorant ClusterManager
// instance. Then, it creates a ClusterManagerCollector for the just created
// ClusterManager. Finally, it registers the ClusterManagerCollector with a
// wrapping Registerer that adds the zone as a label. In this way, the metrics
// collected by different ClusterManagerCollectors do not collide.
func NewUnitStats(reg prometheus.Registerer) *UnitStats {
	c := &UnitStats{}
	sc := UnitStatsCollector{UnitStats: c}
	prometheus.WrapRegistererWith(prometheus.Labels{}, reg).MustRegister(sc)
	return c
}

func main() {
	/*
		var stats *UnitStats
		printResult, err := stats.collectMetrics()
		if err != nil {
			log.Error("Error by main func: ", err.Error())
		}
		fmt.Printf("Response by func: %v\n", printResult)

	*/
	/*
		// For durty debug purposes only
			map[laravel:map[processes:map[idle:1 running:1 starting:0] requests:map[active:0]]
			for application := range printResult.Applications {
				for state := range printResult.Applications[application]["processes"] {
					fmt.Println(printResult.Applications[application]["processes"][state])
				}
				fmt.Println(application)
				fmt.Println(printResult.Applications[application]["processes"])
			}
	*/

	//#############################
	// Since we are dealing with custom Collector implementations, it might
	// be a good idea to try it out with a pedantic registry.
	reg := prometheus.NewPedanticRegistry()

	// Construct cluster managers. In real code, we would assign them to
	// variables to then do something with them.
	NewUnitStats(reg)
	//NewUnitStats(reg)

	// Add the standard process and Go metrics to the custom registry.
	reg.MustRegister(
		prometheus.NewProcessCollector(prometheus.ProcessCollectorOpts{}),
		//prometheus.NewGoCollector(),
	)

	http.Handle("/metrics", promhttp.HandlerFor(reg, promhttp.HandlerOpts{}))
	var metricsAddress = env.GetDefault("METRICS_LISTEN_ADDRESS", ":9095")
	log.Fatal(http.ListenAndServe(metricsAddress, nil))

}

func (stats *UnitStats) collectMetrics() (metrics *UnitStats, err error) {
	//############################
	// Here should be configured with flags or env
	//var network = "unix"
	//var address = "/tmp/unit-sock/control.unit.sock"
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
