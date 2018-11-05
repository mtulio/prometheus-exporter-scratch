package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"sort"

	"github.com/mtulio/prometheus-exporter-scratch/collector"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/common/version"
	log "github.com/sirupsen/logrus"
)

var (
	flagListenAddress *string
	flagMetricsPath   *string
	flagVersion       *bool
	versionStr        string
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func readVersion() {
	dat, err := ioutil.ReadFile("./VERSION")
	check(err)
	versionStr = string(dat)
}

func usage() {
	fmt.Fprintf(os.Stderr, "usage: %s [options]\n", os.Args[0])
	flag.PrintDefaults()
	os.Exit(2)
}

func init() {
	prometheus.MustRegister(version.NewCollector("scratch_exporter"))

	flagListenAddress = flag.String("web.listen-address", ":9999", "Address on which to expose metrics and web interface.")
	flagMetricsPath = flag.String("web.telemetry-path", "/metrics", "Path under which to expose metrics.")
	flagVersion = flag.Bool("v", false, "prints current version")
	flag.Usage = usage
	flag.Parse()

	readVersion()

	if *flagVersion {
		fmt.Println(versionStr)
		os.Exit(0)
	}
}

// Main Prometheus handler
func handler(w http.ResponseWriter, r *http.Request) {
	filters := r.URL.Query()["collect[]"]
	log.Debugln("collect query:", filters)

	mc, err := collector.NewMasterCollector(filters...)
	if err != nil {
		log.Warnln("Couldn't create", err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("Couldn't create %s", err)))
		return
	}

	registry := prometheus.NewRegistry()
	err = registry.Register(mc)
	if err != nil {
		log.Errorln("Couldn't register collector:", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf("Couldn't register collector: %s", err)))
		return
	}

	gatherers := prometheus.Gatherers{
		prometheus.DefaultGatherer,
		registry,
	}
	// Delegate http serving to Prometheus client library, which will call collector.Collect.
	h := promhttp.InstrumentMetricHandler(
		registry,
		promhttp.HandlerFor(gatherers,
			promhttp.HandlerOpts{
				// ErrorLog:      log.NewErrorLogger(),
				ErrorHandling: promhttp.ContinueOnError,
			}),
	)
	h.ServeHTTP(w, r)
}

func main() {

	log.Infoln("Starting exporter ", versionStr)

	// Instance master collector that will keep all subsystems
	// This instance is only used to check collector creation and logging.
	mc, err := collector.NewMasterCollector()
	if err != nil {
		log.Fatalf("Couldn't create collector: %s", err)
	}
	log.Infof("Enabled collectors:")
	collectors := []string{}
	for n := range mc.Collectors {
		collectors = append(collectors, n)
	}
	sort.Strings(collectors)
	for _, n := range collectors {
		log.Infof(" - %s", n)
	}

	//This section will start the HTTP server and expose
	//any metrics on the /metrics endpoint.
	http.HandleFunc(*flagMetricsPath, handler)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`<html>
			<head><title>Prometheus Exporter Scratch</title></head>
			<body>
			<h1>Scratch Exporter</h1>
			<p><a href="` + *flagMetricsPath + `">Metrics</a></p>
			</body>
			</html>`))
	})

	log.Info("Beginning to serve on port " + *flagListenAddress)
	log.Fatal(http.ListenAndServe(*flagListenAddress, nil))
}
