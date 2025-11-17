package main

import (
	_ "github.com/jonasbroms/hbm/storage/driver/sqlite"

	_ "github.com/jonasbroms/hbm/docker/resource/driver/action"
	_ "github.com/jonasbroms/hbm/docker/resource/driver/capability"
	_ "github.com/jonasbroms/hbm/docker/resource/driver/config"
	_ "github.com/jonasbroms/hbm/docker/resource/driver/device"
	_ "github.com/jonasbroms/hbm/docker/resource/driver/dns"
	_ "github.com/jonasbroms/hbm/docker/resource/driver/image"
	_ "github.com/jonasbroms/hbm/docker/resource/driver/logdriver"
	_ "github.com/jonasbroms/hbm/docker/resource/driver/logopt"
	_ "github.com/jonasbroms/hbm/docker/resource/driver/plugin"
	_ "github.com/jonasbroms/hbm/docker/resource/driver/port"
	_ "github.com/jonasbroms/hbm/docker/resource/driver/registry"
	_ "github.com/jonasbroms/hbm/docker/resource/driver/volume"
	_ "github.com/jonasbroms/hbm/docker/resource/driver/volumedriver"
)
