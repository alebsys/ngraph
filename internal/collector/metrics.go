package collector

import (
	"fmt"
	"strings"
)

const (
	inputMetricName  = "network_connections_input_total"
	outputMetricName = "network_connections_output_total"
)

func createMetric(k string, v int) (string, error) {
	connMeta := strings.Split(k, "-")
	if connMeta[2] == "output" {
		return fmt.Sprintf("%s{src_ip=\"%s\", dest_ip=\"%s\"} %d\n", outputMetricName, connMeta[0], connMeta[1], v), nil
	}
	return fmt.Sprintf("%s{src_ip=\"%s\", dest_ip=\"%s\"} %d\n", inputMetricName, connMeta[1], connMeta[0], v), nil
}
