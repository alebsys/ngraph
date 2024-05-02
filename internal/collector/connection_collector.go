package collector

import (
	"fmt"
	"strings"
	"syscall"

	lnet "github.com/alebsys/ngraph/internal/localnetinfo"
	"github.com/prometheus/procfs"
	"github.com/vishvananda/netlink"
	"github.com/vishvananda/netns"
)

const (
	IpLocalPortRangeFile = "/proc/sys/net/ipv4/ip_local_port_range"
)

// getNsPids return all pids owner network namespaces
func (c *Collector) getConnections() (map[string]int, error) {
	fs, err := procfs.NewFS("/proc")
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
	// networkNamespacePIDs := make(map[uint32]int)
	connections := make(map[string]int)

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

// getNetworkNamespaceID retrieves the inode of the network namespace associated with a given process.
func (c *Collector) getNetworkNamespaceID(process procfs.Proc) (uint32, error) {
	namespaces, err := process.Namespaces()
	if err != nil {
		return 0, err
	}
	return selectNetworkNamespaceInode(namespaces), nil
}

// selectNetworkNamespaceInode selects the inode of the network namespace from the list of namespaces.
func selectNetworkNamespaceInode(namespaces procfs.Namespaces) uint32 {
	for _, namespace := range namespaces {
		// we are only interested in the inode of the network namespace
		if namespace.Type == "net" {
			return namespace.Inode
		}
	}
	return 0
}

func (c *Collector) getEstabConnectionsFromNetNs(portRanges lnet.LocalPortRange, connections *map[string]int, ns netns.NsHandle) error {
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
		// Check if the connection is established (St == 1)
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

		direction := checkConnectDirection(int(conn.InetDiagMsg.ID.SourcePort), portRanges.MinPort, portRanges.MaxPort)

		key := fmt.Sprintf("%s-%s-%s", conn.InetDiagMsg.ID.Source, conn.InetDiagMsg.ID.Destination, direction)
		(*connections)[key]++
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
