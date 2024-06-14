package localnetinfo

import (
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
)

type LocalPortRange struct {
	MinPort int
	MaxPort int
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

func GetPortRange(file string) (LocalPortRange, error) {
	data, err := os.ReadFile(file)
	if err != nil {
		return LocalPortRange{}, err
	}

	ports := strings.Fields(string(data))

	min, err := strconv.Atoi(ports[0])
	if err != nil {
		return LocalPortRange{}, err
	}

	max, err := strconv.Atoi(ports[1])
	if err != nil {
		return LocalPortRange{}, err
	}

	return LocalPortRange{
		MinPort: min,
		MaxPort: max,
	}, nil
}
