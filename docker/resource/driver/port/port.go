package port

import (
	"fmt"

	"github.com/jonasbroms/hbm/docker/resource"
	"github.com/jonasbroms/hbm/docker/resource/driver"
	"github.com/jonasbroms/hbm/pkg/utils"
)

type Config struct{}

func init() {
	resource.RegisterDriver("port", New)
}

func New() (driver.Resourcer, error) {
	return &Config{}, nil
}

func (c *Config) List() interface{} {
	return []string{}
}

func (c *Config) Valid(value string) error {
	var ports []int

	startPort, endPort, err := utils.GetPortRangeFromString(value)
	if err != nil {
		return err
	}

	ports = append(ports, startPort)

	if startPort != endPort {
		if startPort > endPort {
			return fmt.Errorf("Range of ports is not valid. Start port is greater than end port.")
		}

		ports = append(ports, endPort)
	}

	for _, p := range ports {
		if p < 0 || p > 65535 {
			return fmt.Errorf("port is not valid")
		}
	}

	return nil
}

func (c *Config) ValidOptions(options map[string]string) error {
	return nil
}
