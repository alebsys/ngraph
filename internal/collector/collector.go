package collector

import (
	"fmt"
	"net"
	"strings"
	"syscall"

	lnet "github.com/alebsys/ngraph/internal/localnetinfo"
	"github.com/prometheus/procfs"
	"github.com/vishvananda/netlink"
	"github.com/vishvananda/netns"
)

// Define private IP ranges as a package-level variable
var privateIPRanges = []string{
	"127.0.0.0/8",
	"10.0.0.0/8",
	"172.16.0.0/12",
	"192.168.0.0/16",
	"100.64.0.0/10",  // Carrier-grade NAT
	"169.254.0.0/16", // Link-local addresses
}

const (
	IpLocalPortRangeFile = "/proc/sys/net/ipv4/ip_local_port_range"
)

type Collector struct {
	cfg Config
}

type Config struct {
	ConnectFromAllNs bool
	ExcludeSubnets   []string
	AllowPublicIP    bool
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
func NewConfig(exclude string, allNS, allowPub bool) *Config {
	excludeSubnets := strings.Split(exclude, ",")

	if len(excludeSubnets) == 1 && excludeSubnets[0] == "" {
		excludeSubnets[0] = "none"
	}
	return &Config{
		ConnectFromAllNs: allNS,
		ExcludeSubnets:   excludeSubnets,
		AllowPublicIP:    allowPub,
	}
}

// GetConnections() TODO:
func (c *Collector) GetConnections() (map[UniqTupleConnection]float64, error) {
	procs, err := procfs.AllProcs()
	if err != nil {
		return nil, err
	}

	portRanges, err := lnet.GetPortRange(IpLocalPortRangeFile)
	if err != nil {
		return nil, err
	}

	localIP, err := lnet.GetLocalIP()
	if err != nil {
		return nil, err
	}

	// key - inode network namespace, value - PID owner
	networkNamespacePIDs := make(map[string]int)
	connections := make(map[UniqTupleConnection]float64)

	for _, proc := range procs {
		ns, err := netns.GetFromPid(proc.PID)
		if err != nil {
			continue
		}
		defer ns.Close()

		// Skip if the network namespace is already processed
		if _, ok := networkNamespacePIDs[ns.UniqueId()]; ok {
			continue
		}
		networkNamespacePIDs[ns.UniqueId()] = proc.PID

		if err := c.getEstabConnectionsFromNetNs(portRanges, &connections, ns, localIP); err != nil {
			continue
		}

		// If ConnectFromAllNs == false (from config) then get connections only from root namespace (PID 1)
		if proc.PID == 1 && !c.cfg.ConnectFromAllNs {
			break
		}
	}
	return connections, nil
}

// getEstabConnectionsFromNetNs TODO:
func (c *Collector) getEstabConnectionsFromNetNs(portRanges lnet.LocalPortRange, connections *map[UniqTupleConnection]float64, ns netns.NsHandle, localIP string) error {
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
		peerIP := conn.InetDiagMsg.ID.Destination.String()

		uniqConn.Direction = checkConnectDirection(int(conn.InetDiagMsg.ID.SourcePort), portRanges.MinPort, portRanges.MaxPort)
		uniqConn.SrcIP = localIP
		uniqConn.DstIP, err = determineDstIP(peerIP, c.cfg.AllowPublicIP)
		if err != nil {
			return err
		}
		(*connections)[uniqConn]++
	}
	return nil
}

// determineDstIP returns the peer IP if public IPs are allowed or if the peer IP is private; otherwise, it returns "external_ip".
func determineDstIP(peerIP string, allowPublicIP bool) (string, error) {
	if allowPublicIP {
		return peerIP, nil
	}

	isPrivate, err := isPrivateIP(peerIP)
	if err != nil {
		return "", err
	}

	if isPrivate {
		return peerIP, nil
	}
	return "external_ip", nil
}

// isPrivateIP checks if the given IP address is from a private range.
func isPrivateIP(ipStr string) (bool, error) {
	ip := net.ParseIP(ipStr)
	if ip == nil {
		return false, fmt.Errorf("invalid IP address: %s", ipStr)
	}

	for _, cidr := range privateIPRanges {
		_, subnet, err := net.ParseCIDR(cidr)
		if err != nil {
			return false, err
		}
		if subnet.Contains(ip) {
			return true, nil
		}
	}

	return false, nil
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
