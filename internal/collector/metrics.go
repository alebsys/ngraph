package collector

import (
	"fmt"
	"log"
	"strings"

	"github.com/alebsys/ngraph/internal/utils"
)

const (
	inputMetricName  = "network_connections_input_total"
	outputMetricName = "network_connections_output_total"
)

func createMetric(k string, v int) (string, error) {
	connMeta := strings.Split(k, "-")
	peerHostName, err := utils.ResolveAddr(connMeta[1])
	if err != nil {
		log.Printf("error: the address %s was not resolved", connMeta[1])
	}

	if connMeta[2] == "output" {
		return fmt.Sprintf("%s{src_ip=\"%s\", dest_ip=\"%s\", peer_hostname=\"%s\"} %d\n", outputMetricName, connMeta[0], connMeta[1], peerHostName, v), nil
	}
	return fmt.Sprintf("%s{src_ip=\"%s\", dest_ip=\"%s\", peer_hostname=\"%s\"} %d\n", inputMetricName, connMeta[1], connMeta[0], peerHostName, v), nil
}
