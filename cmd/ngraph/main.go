package main

import (
	"log"

	"github.com/alebsys/ngraph/internal/collector"
	"github.com/alecthomas/kingpin/v2"
)

func main() {
	log.Println("ngraph is starting")

	scrapeInterval := kingpin.Flag("interval", "Interval for generating metrics -- ONE SHOT version").Default("15").Int()
	metricsFilePath := kingpin.Flag("output", "Path for metric files").Required().String()
	connectFromAllNs := kingpin.Flag("all", "Scrape connections from all network namespaces").Default("false").Bool()
	excludeSubnets := kingpin.Flag("exclude", "Comma separated list of pattern subnets to skip them during connection parsing, example: --exclude=127.0,192.168").Default("none").String()
	kingpin.HelpFlag.Short('h')
	kingpin.Parse()

	config := collector.NewConfig(*scrapeInterval, *metricsFilePath, *excludeSubnets, *connectFromAllNs)
	c := collector.NewCollector(config)
	c.Run()
}
