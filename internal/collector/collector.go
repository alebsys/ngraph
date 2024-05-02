package collector

import (
	"strings"
	"syscall"

	lnet "github.com/alebsys/ngraph/internal/localnetinfo"
	"github.com/prometheus/procfs"
	"github.com/vishvananda/netlink"
	"github.com/vishvananda/netns"
)

const (
	IpLocalPortRangeFile = "/proc/sys/net/ipv4/ip_local_port_range"
	procRoot             = "/proc"
)

type Collector struct {
	cfg Config
}

type Config struct {
	ConnectFromAllNs bool
	ExcludeSubnets   []string
}

type UniqTupleConnection struct {
	SrcIP     string
	DstIP     string
	Direction string
}

// NewCollector creates a new Collector instance with the given configuration.
func NewCollector(c Config) *Collector {
	return &Collector{
		cfg: c,
	}
}

// NewConfig creates a new Config instance for Collector with the given fields.
func NewConfig(exclude string, allNS bool) *Config {
	excludeSubnets := strings.Split(exclude, ",")

	if len(excludeSubnets) == 1 && excludeSubnets[0] == "" {
		excludeSubnets[0] = "none"
	}
	return &Config{
		ConnectFromAllNs: allNS,
		ExcludeSubnets:   excludeSubnets,
	}
}

// GetConnections() TODO:
func (c *Collector) GetConnections() (map[UniqTupleConnection]float64, error) {
	fs, err := procfs.NewFS(procRoot)
	if err != nil {
		return nil, err
	}

	// get all processes
	processes, err := fs.AllProcs()
	if err != nil {
		return nil, err
	}

	portRanges, err := lnet.GetPortRange(IpLocalPortRangeFile)
	if err != nil {
		return nil, err
	}

	// key - inode network namespace, value - PID owner
	networkNamespacePIDs := make(map[string]int)
	connections := make(map[UniqTupleConnection]float64)

	for _, process := range processes {
		ns, err := netns.GetFromPid(process.PID)
		if err != nil {
			continue
		}

		// Skip if the network namespace is already processed
		if _, ok := networkNamespacePIDs[ns.UniqueId()]; ok {
			continue
		}
		networkNamespacePIDs[ns.UniqueId()] = process.PID

		if err := c.getEstabConnectionsFromNetNs(portRanges, &connections, ns); err != nil {
			continue
		}

		// If ConnectFromAllNs == false (config) get connections only from root namespace (PID 1)
		if process.PID == 1 && !c.cfg.ConnectFromAllNs {
			break
		}
	}
	return connections, nil
}

// getEstabConnectionsFromNetNs TODO:
func (c *Collector) getEstabConnectionsFromNetNs(portRanges lnet.LocalPortRange, connections *map[UniqTupleConnection]float64, ns netns.NsHandle) error {
	err := netns.Setns(ns, syscall.CLONE_NEWNET)
	if err != nil {
		return err
	}

	var conns []*netlink.InetDiagTCPInfoResp

	ipv4Conns, err := netlink.SocketDiagTCPInfo(syscall.AF_INET)
	if err != nil {
		return err
	}

	ipv6Conns, err := netlink.SocketDiagTCPInfo(syscall.AF_INET6)
	if err != nil {
		return err
	}

	conns = append(conns, ipv4Conns...)
	conns = append(conns, ipv6Conns...)

	for _, conn := range conns {
		uniqConn := UniqTupleConnection{}
		if conn.InetDiagMsg.State != 1 {
			continue
		}
		if conn.TCPInfo == nil {
			continue
		}

		// Check if the connection addresses are within the exclude subnets to ignore it
		if shouldMatchBySubnets(c.cfg.ExcludeSubnets, conn.InetDiagMsg.ID.Source.String()) || shouldMatchBySubnets(c.cfg.ExcludeSubnets, conn.InetDiagMsg.ID.Destination.String()) {
			continue
		}

		uniqConn.Direction = checkConnectDirection(int(conn.InetDiagMsg.ID.SourcePort), portRanges.MinPort, portRanges.MaxPort)
		uniqConn.SrcIP = conn.InetDiagMsg.ID.Source.String()
		uniqConn.DstIP = conn.InetDiagMsg.ID.Destination.String()
		(*connections)[uniqConn]++
	}
	return nil
}

func shouldMatchBySubnets(subnets []string, addr string) bool {
	for _, s := range subnets {
		// handle cases: --exclude ("10.32.68,") || (",10.32.68,")
		if len(s) == 0 {
			continue
		}
		if strings.HasPrefix(addr, s) {
			return true
		}
	}
	return false
}

// checkConnectDirection determines the direction of a network connection based on the specified port range.
// If the given port is within the allowed range, the function returns "output" indicating an outgoing connection.
// If the port is outside the allowed range, the function returns "input" indicating an incoming connection.
func checkConnectDirection(p, min, max int) string {
	if p >= min && p <= max {
		return "output"
	}
	return "input"
}
