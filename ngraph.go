package main

import (
	"log"
	"net/http"

	"github.com/alebsys/ngraph/internal/collector"
	"github.com/alebsys/ngraph/internal/exporter"
	"github.com/alecthomas/kingpin/v2"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	listeningAddress := kingpin.Flag("address", "Address on which to expose metrics.").Default(":9234").String()
	metricsEndpoint := kingpin.Flag("endpoint", "Path under which to expose metrics.").Default("/metrics").String()
	excludeSubnets := kingpin.Flag("exclude", "Comma separated list of pattern subnets to skip them during connection parsing, example: --exclude=127.0,192.168").Default("none").String()
	connectFromAllNs := kingpin.Flag("all", "Scrape connections from all network namespaces").Default("false").Bool()
	kingpin.HelpFlag.Short('h')
	kingpin.Parse()

	log.Printf("ngraph is starting on %s", *listeningAddress)

	registerer := prometheus.DefaultRegisterer
	gatherer := prometheus.DefaultGatherer

	collCfg := collector.NewConfig(*excludeSubnets, *connectFromAllNs)
	coll := collector.NewCollector(*collCfg)

	exporterCfg := exporter.NewConfig(*listeningAddress, *metricsEndpoint)
	exporter := exporter.NewExporter(*exporterCfg, coll)

	registerer.MustRegister(exporter)
	http.Handle(*metricsEndpoint, promhttp.HandlerFor(gatherer, promhttp.HandlerOpts{}))
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`<html>
				<head><title>ngraph exporter</title></head>
				<body>
				<h1>ngraph exporter</h1>
				<p><a href="` + *metricsEndpoint + `">Metrics</a></p>
				</body>
				</html>`))
	})
	log.Fatal(http.ListenAndServe(*listeningAddress, nil))
}
