package main

import (
	"flag"
	"log"

	"github.com/alebsys/ngraph/internal/collector"
)

var (
	config = &collector.Config{}
)

func init() {
	flag.IntVar(&config.ScrapeInterval, "interval", collector.DefaultScrapeInterval, "Interval for generating metrics")
	flag.StringVar(&config.MetricsFilePath, "output", collector.DefaultMetricsFilePath, "Path for metric files")
	flag.BoolVar(&config.ConnectFromAllNs, "all", false, "Scrape connections from all network namespaces")
	flag.Parse()
}
func main() {
	log.Println("ngraph is starting")

	c := collector.NewCollector(config)
	c.Run()
}
