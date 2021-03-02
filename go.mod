module github.com/kassisol/hbm

go 1.16

require (
	github.com/Microsoft/go-winio v0.4.8-0.20180501170546-ab35fc04b636 // indirect
	github.com/coreos/go-systemd v0.0.0-20180409111510-d1b7d058aa2a // indirect
	github.com/cpuguy83/go-md2man v1.0.8 // indirect
	github.com/docker/distribution v2.6.0-rc.1.0.20180327202408-83389a148052+incompatible // indirect
	github.com/docker/docker v17.12.0-ce-rc1.0.20180507131252-57493cd60628+incompatible
	github.com/docker/go-connections v0.4.0
	github.com/docker/go-plugins-helpers v0.0.0-20180116160015-61cb8e233420
	github.com/docker/go-units v0.3.3 // indirect
	github.com/gogo/protobuf v1.0.1-0.20180330174643-1ef32a8b9fc3 // indirect
	github.com/google/go-cmp v0.5.4 // indirect
	github.com/gotestyourself/gotestyourself v2.2.0+incompatible // indirect
	github.com/inconshreveable/mousetrap v1.0.0 // indirect
	github.com/jinzhu/gorm v1.9.17-0.20200921022817-466b344ff592
	github.com/jinzhu/inflection v1.0.1-0.20210111022912-b5281034e75e // indirect
	github.com/juliengk/go-docker v0.0.0-20180321163138-f5cda316edb7
	github.com/juliengk/go-log v0.0.0-20171002012451-57b916563bbb
	github.com/juliengk/go-mount v0.0.0-20170406141235-e7123cbaaaf6
	github.com/juliengk/go-utils v0.0.0-00010101000000-000000000000
	github.com/juliengk/go-utils/validation v0.0.0-00010101000000-000000000000
	github.com/mattn/go-sqlite3 v1.14.7-0.20210218183441-ab91e9342b02 // indirect
	github.com/opencontainers/go-digest v1.0.0-rc1.0.20180430190053-c9281466c8b2 // indirect
	github.com/opencontainers/image-spec v1.0.1-0.20180411145040-e562b0440392 // indirect
	github.com/pkg/errors v0.8.1-0.20180311214515-816c9085562c // indirect
	github.com/russross/blackfriday v1.5.2-0.20180428102519-11635eb403ff // indirect
	github.com/sirupsen/logrus v1.0.6-0.20180329225952-778f2e774c72
	github.com/spf13/cobra v0.0.3
	github.com/spf13/pflag v1.0.1 // indirect
	github.com/stretchr/testify v1.7.0 // indirect
	gopkg.in/airbrake/gobrake.v2 v2.0.9 // indirect
	gopkg.in/gemnasium/logrus-airbrake-hook.v2 v2.1.2 // indirect
	gopkg.in/yaml.v2 v2.2.1 // indirect
	gotest.tools v2.2.0+incompatible // indirect
)

replace github.com/juliengk/go-utils => ./local/github.com/juliengk/go-utils/

replace github.com/juliengk/go-utils/filedir => ./local/github.com/juliengk/go-utils/filedir

replace github.com/juliengk/go-utils/json => ./local/github.com/juliengk/go-utils/json

replace github.com/juliengk/go-utils/validation => ./local/github.com/juliengk/go-utils/validation
