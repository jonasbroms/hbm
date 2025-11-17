package image

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
	resource.RegisterDriver("image", New)
}

func New() (driver.Resourcer, error) {
	keys := []string{"subimages"}

	return &Config{Options: keys}, nil
}

func (c *Config) List() interface{} {
	return []string{}
}

func (c *Config) Valid(value string) error {
	return nil
}

func (c *Config) ValidOptions(options map[string]string) error {

	for k := range options {
		if !utils.StringInSlice(k, c.Options, false) {
			return fmt.Errorf("%s is not a valid option key", k)
		}
	}
	return nil
}
