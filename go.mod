module github.com/kassisol/hbm

go 1.23.0

toolchain go1.23.6

require (
	github.com/docker/docker v27.3.1+incompatible
	github.com/docker/go-connections v0.5.0
	github.com/docker/go-plugins-helpers v0.0.0-20211224144127-6eecb7beb651
	github.com/jinzhu/gorm v1.9.17-0.20211120011537-5c235b72a414
	github.com/juliengk/go-docker v0.0.0-20180321163138-f5cda316edb7
	github.com/juliengk/go-mount v0.0.0-20170406141235-e7123cbaaaf6
	github.com/juliengk/go-utils v0.0.0-00010101000000-000000000000
	github.com/sirupsen/logrus v1.9.3
	github.com/spf13/cobra v1.8.0
)

require (
	github.com/Microsoft/go-winio v0.6.1 // indirect
	github.com/containerd/log v0.1.0 // indirect
	github.com/coreos/go-systemd v0.0.0-20180409111510-d1b7d058aa2a // indirect
	github.com/cpuguy83/go-md2man/v2 v2.0.3 // indirect
	github.com/distribution/reference v0.5.0 // indirect
	github.com/docker/go-units v0.5.0 // indirect
	github.com/felixge/httpsnoop v1.0.4 // indirect
	github.com/go-logr/logr v1.4.2 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/inconshreveable/mousetrap v1.1.0 // indirect
	github.com/jinzhu/inflection v1.0.1-0.20231016084002-bbe0a3e7399f // indirect
	github.com/mattn/go-sqlite3 v1.14.19 // indirect
	github.com/moby/docker-image-spec v1.3.1 // indirect
	github.com/moby/term v0.5.0 // indirect
	github.com/morikuni/aec v1.0.0 // indirect
	github.com/opencontainers/go-digest v1.0.0 // indirect
	github.com/opencontainers/image-spec v1.0.2 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/russross/blackfriday/v2 v2.1.0 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	go.opentelemetry.io/auto/sdk v1.1.0 // indirect
	go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp v0.60.0 // indirect
	go.opentelemetry.io/otel v1.35.0 // indirect
	go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp v1.35.0 // indirect
	go.opentelemetry.io/otel/metric v1.35.0 // indirect
	go.opentelemetry.io/otel/trace v1.35.0 // indirect
	golang.org/x/mod v0.17.0 // indirect
	golang.org/x/sync v0.11.0 // indirect
	golang.org/x/sys v0.30.0 // indirect
	golang.org/x/text v0.22.0 // indirect
	golang.org/x/time v0.5.0 // indirect
	golang.org/x/tools v0.21.1-0.20240508182429-e35e4ccd0d2d // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
	gotest.tools/v3 v3.5.1 // indirect
)

replace github.com/juliengk/go-utils => github.com/jonasbroms/go-utils v0.0.0-20220910194142-e74a4beb9866

replace github.com/juliengk/go-docker => github.com/kassisol/go-docker v0.0.0-20180321163138-f5cda316edb7

replace github.com/juliengk/go-mount => github.com/kassisol/go-mount v0.0.0-20170406141235-e7123cbaaaf6
