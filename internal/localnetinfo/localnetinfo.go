package localnetinfo

import (
	"os"
	"strconv"
	"strings"
)

type LocalPortRange struct {
	MinPort int
	MaxPort int
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
