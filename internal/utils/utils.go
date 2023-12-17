package utils

import (
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
	return hostnames[0], nil
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

// TODO:
func GetMainHostIpAddress() (string, error) {
	h, err := os.Hostname()
	if err != nil {
		return "", err
	}

	ips, err := net.LookupIP(h)
	if err != nil {
		return "", err
	}
	// TODO: выбираем всегда первый элемент?
	return ips[0].String(), nil
}
