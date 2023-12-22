package localnetinfo

import (
	"os"
	"strconv"
	"strings"
)

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
