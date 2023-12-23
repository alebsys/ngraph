package collector

import (
	"fmt"
	"strings"

	lnet "github.com/alebsys/ngraph/internal/localnetinfo"
	"github.com/prometheus/procfs"
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
	networkNamespacePIDs := make(map[uint32]int)
	connections := make(map[string]int)

	for _, process := range processes {
		netNsID, err := c.getNetworkNamespaceID(process)
		if err != nil || netNsID == 0 {
			continue
		}

		// Skip if the network namespace is already processed
		if _, ok := networkNamespacePIDs[netNsID]; ok {
			continue
		}
		networkNamespacePIDs[netNsID] = process.PID

		if err := c.getConnectionsFromNamespace(portRanges, &connections, process); err != nil {
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

// getConnectionsFromNamespace retrieves and counts established connections from a network namespace.
// It filters connections based on the established status and the specified port range.
func (c *Collector) getConnectionsFromNamespace(portRanges lnet.LocalPortRange, connections *map[string]int, process procfs.Proc) error {
	// Get all connections from /proc/<pid>/net/tcp (per network namespace)
	conns, err := process.NetTCP()
	if err != nil {
		return err
	}

	for _, conn := range conns.NetTCP {
		// Check if the connection is established (St == 1)
		if conn.St != 1 {
			continue
		}
		// Check if the connection addresses are within the exclude subnets to ignore it
		if shouldMatchBySubnets(c.cfg.ExcludeSubnets, conn.LocalAddr.String()) || shouldMatchBySubnets(c.cfg.ExcludeSubnets, conn.RemAddr.String()) {
			continue
		}
		connDirection := checkConnectDirection(int(conn.LocalPort), portRanges.MinPort, portRanges.MaxPort)

		key := fmt.Sprintf("%s-%s-%s", conn.LocalAddr, conn.RemAddr, connDirection)
		(*connections)[key]++
	}
	return nil
}

func shouldMatchBySubnets(subnets []string, addr string) bool {
	for _, s := range subnets {
		// // handle cases: --exclude ("10.32.68,") || (",10.32.68,")
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
