package volume

import (
	"fmt"

	"github.com/jonasbroms/hbm/docker/resource"
	"github.com/jonasbroms/hbm/docker/resource/driver"
	"github.com/juliengk/go-utils"
)

type Config struct {
	Options []string
}

func init() {
	resource.RegisterDriver("volume", New)
}

func New() (driver.Resourcer, error) {
	keys := []string{
		"recursive",
		"nosuid",
	}

	return &Config{Options: keys}, nil
}

func (c *Config) List() interface{} {
	return []string{}
}

func (c *Config) Valid(value string) error {
	return nil
}

func (c *Config) ValidOptions(options map[string]string) error {
	if len(options) == 0 {
		return nil
	}

	for k := range options {
		if !utils.StringInSlice(k, c.Options, false) {
			return fmt.Errorf("%s is not a valid option key", k)
			//fmt.Printf("Conflicting options --type %s and --recursive\n", resourceAddType)
		}
	}

	return nil
}
