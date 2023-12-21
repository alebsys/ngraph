package localnetinfo

import (
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
)

// TODO: проверить какое значение из slice следует возвращать.
func ResolveAddr(addr string) (string, error) {
	hostnames, err := net.LookupAddr(addr)
	if err != nil {
		return "", err
	}
	hostname := hostnames[0]

	// return hostname without last dot symbol
	if len(hostname) > 0 {
		hostname = hostname[:len(hostnames)-1]
	}
	return hostname, nil
}

func GetPortRange(file string) (int, int, error) {
	data, err := os.ReadFile(file)
	if err != nil {
		return 0, 0, err
	}

	ports := strings.Fields(string(data))

	min, err := strconv.Atoi(ports[0])
	if err != nil {
		return 0, 0, err
	}

	max, err := strconv.Atoi(ports[1])
	if err != nil {
		return 0, 0, err
	}

	return min, max, nil
}

// GetLocalIP returns the non loopback local IP of the host
func GetLocalIP() (string, error) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return "", fmt.Errorf("failed to get local IP: %v", err)
	}
	for _, a := range addrs {
		// check the address type and if it is not a loopback the display it
		if ipnet, ok := a.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String(), nil
			}
		}
	}
	return "", fmt.Errorf("local IP not found")
}
