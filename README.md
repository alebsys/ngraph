# ngraph - Active TCP Connections Collector

`ngraph` is a metrics collector tool designed to work in conjunction with [node_exporter](https://github.com/prometheus/node_exporter). Its primary goal is to provide information about existing network connections on a machine, which can then be utilized to build a graph of network interactions within your infrastructure.

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

- `--interval`: Interval for generating metrics (default: 15 seconds).
- `--output`: Path for metric files (default: "/tmp/node_exporter/textfile").
- `--all`: Scrape connections from all network namespaces (default: false).

Example:

```bash
./ngraph --interval=30 --output=/path/to/metrics --all
```

## Metrics

`ngraph` generates Prometheus-compatible metrics in a textfile format. The resulting metrics are stored in the specified output directory with the filename `ngraph.prom`.

The primary metrics provided is:

- `network_connections_input_total`: Total input number of unique network connections per hosts;
- `network_connections_output_total`: Total output number of unique network connections per hosts.

Example:

```bash
network_connections_input_total{dest_ip="10.12.57.104", src_ip="10.12.57.27", peer_hostname="example.com"} 3
network_connections_output_total{src_ip="10.12.57.104", dest_ip="10.24.127.27", peer_hostname="example.com"} 2
```

Example of possible visualization of metrics in grafana:

![Example of possible visualization of metrics in grafana](https://file.notion.so/f/f/2c3f117c-ca9f-4cb1-af7b-a51fc6db39bb/ebed4d5a-550f-4a55-8fd8-26df7f27a418/Untitled.png?id=55aa7c53-db6a-4152-a992-ad60bec209ba&table=block&spaceId=2c3f117c-ca9f-4cb1-af7b-a51fc6db39bb&expirationTimestamp=1703030400000&signature=yWncq7i-U8-ITGAku_QMBpt0Us8R5EAOAT3Ggh9bmvw&downloadName=Untitled.png)

## Contributing

Feel free to contribute to `ngraph` by opening issues, proposing new features, or submitting pull requests. Your feedback and contributions are highly valued.

## License

This project is licensed under the [MIT License](LICENSE).

## Acknowledgments

`ngraph` makes use of the [procfs](https://github.com/prometheus/procfs) library to interact with the `/proc` filesystem. Special thanks to the authors and contributors of `procfs`.
