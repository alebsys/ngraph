package collector

import (
	"fmt"

	"github.com/alebsys/ngraph/internal/utils"
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

	minPort, maxPort, err := utils.GetPortRange(IpLocalPortRangeFile)
	if err != nil {
		return nil, err
	}

	// TODO: добавить описание к функции и в целом зачем нужен ip адрес хоста
	localIP, err := utils.GetLocalIP()
	if err != nil {
		return nil, err
	}

	// key - inode network namespace, value - PID owner
	networkNamespacePIDs := make(map[uint32]int)
	connections := make(map[string]int)

	for _, process := range processes {
		netNsID, err := c.getNetworkNamespaceID(process)
		if err != nil {
			return nil, err
		}

		// Skip if the network namespace is already processed
		if _, ok := networkNamespacePIDs[netNsID]; ok {
			continue
		}
		networkNamespacePIDs[netNsID] = process.PID

		if err := c.getConnectionsFromNamespace(minPort, maxPort, &connections, process, localIP); err != nil {
			return nil, err
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
	for _, namespace := range namespaces {
		// we are only interested in the inode of the network namespace
		if namespace.Type == "net" {
			return namespace.Inode, nil
		}
	}
	return 0, nil
}

// getConnectionsFromNamespace retrieves and counts established connections from a network namespace.
// It filters connections based on the established status and the specified port range.
func (c *Collector) getConnectionsFromNamespace(minPort, maxPort int, connections *map[string]int, process procfs.Proc, localIP string) error {
	// Get all connections from /proc/<pid>/net/tcp (per network namespace)
	conns, err := process.NetTCP()
	if err != nil {
		return err
	}

	for _, tcpConns := range conns.NetTCP {
		// Check if the connection is established (St == 1)
		if tcpConns.St != 1 {
			continue
		}
		connectDirection := checkConnectDirection(int(tcpConns.LocalPort), minPort, maxPort)

		key := fmt.Sprintf("%s-%s-%s", localIP, tcpConns.RemAddr, connectDirection)
		(*connections)[key]++
	}
	return nil
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
