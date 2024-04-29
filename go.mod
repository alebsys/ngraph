module github.com/alebsys/ngraph

go 1.19

require (
	github.com/prometheus/client_golang v1.19.0
	github.com/prometheus/procfs v0.12.0
)

require (
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/cespare/xxhash/v2 v2.2.0 // indirect
	github.com/prometheus/client_model v0.5.0 // indirect
	github.com/prometheus/common v0.48.0 // indirect
	github.com/vishvananda/netns v0.0.0-20200728191858-db3c7e526aae
	google.golang.org/protobuf v1.32.0 // indirect
)

require (
	github.com/alecthomas/kingpin/v2 v2.4.0
	github.com/alecthomas/units v0.0.0-20211218093645-b94a6e3cc137 // indirect
	github.com/vishvananda/netlink v1.2.1-beta.2
	github.com/xhit/go-str2duration/v2 v2.1.0 // indirect
	golang.org/x/sys v0.16.0 // indirect
)

replace github.com/prometheus/procfs v0.12.0 => github.com/alebsys/procfs v0.12.9
