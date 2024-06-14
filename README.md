# ngraph - Active TCP Connections Collector

`ngraph` is a metrics exporter. Its primary goal is to provide information about existing network connections on a machine, which can then be utilized to build a graph of network interactions within your infrastructure.

## Installation

1. Clone the repository:

   ```bash
   git clone https://github.com/alebsys/ngraph.git
   ```

2. Build the `ngraph` binary:

   ```bash
   cd ngraph
   go build
   ```

3. Run the `ngraph` binary:

   ```bash
   ./ngraph
   ```

## Usage

`ngraph` accepts command-line arguments to customize its behavior. Below are the available options:

```
./ngraph --help
usage: ngraph [<flags>]


Flags:
  -h, --[no-]help            Show context-sensitive help (also try --help-long and --help-man).
      --address=":9234"      Address on which to expose metrics.
      --endpoint="/metrics"  Path under which to expose metrics.
      --exclude="none"       Comma separated list of pattern subnets to skip them during connection parsing, example: --exclude=127.0,192.168
      --[no-]all             Scrape connections from all network namespaces
```

Example:

```bash
./ngraph --all
```

## Metrics

`ngraph` generates Prometheus-compatible metrics in a textfile format. The resulting metrics are stored in the specified output directory with the filename `ngraph.prom`.

The primary metrics provided is:

- `network_connections_incoming_total`: Total input number of unique network connections per hosts;
- `network_connections_outgoing_total`: Total output number of unique network connections per hosts.

Example:

```bash
# HELP network_connections_incoming_total Total number of incoming network connections between source and destination IP addresses.
# TYPE network_connections_incoming_total gauge
network_connections_incoming_total{dst_ip="10.13.77.11",src_ip="10.13.70.153"} 2
network_connections_incoming_total{dst_ip="10.13.77.11",src_ip="10.13.71.54"} 1
# HELP network_connections_outgoing_total Total number of outgoing network connections between source and destination IP addresses.
# TYPE network_connections_outgoing_total gauge
network_connections_outgoing_total{dst_ip="10.13.24.3",src_ip="10.13.77.11"} 1
```

## Contributing

Feel free to contribute to `ngraph` by opening issues, proposing new features, or submitting pull requests. Your feedback and contributions are highly valued.

## License

This project is licensed under the [MIT License](LICENSE).

## Acknowledgments

`ngraph` makes use of the [procfs](https://github.com/prometheus/procfs) library to interact with the `/proc` filesystem. Special thanks to the authors and contributors of `procfs`.
