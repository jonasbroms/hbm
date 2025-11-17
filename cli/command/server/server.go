package server

import (
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"path"
	"syscall"

	"github.com/docker/go-plugins-helpers/authorization"
	"github.com/jonasbroms/hbm/pkg/adf"
	"github.com/jonasbroms/hbm/plugin"
	"github.com/jonasbroms/hbm/version"
	"github.com/juliengk/go-utils/filedir"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func NewServerCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "server",
		Short: "Launch the HBM server",
		Long:  serverDescription,
		Args:  cobra.NoArgs,
		Run:   runStart,
	}

	return cmd
}

func serverInitConfig() {
	dockerPluginPath := "/etc/docker/plugins"
	dockerPluginFile := path.Join(dockerPluginPath, "hbm.spec")
	pluginSpecContent := []byte("unix://run/docker/plugins/hbm.sock")

	_, err := exec.LookPath("docker")
	if err != nil {
		fmt.Println("Docker does not seem to be installed. Please check your installation.")

		os.Exit(-1)
	}

	if err := filedir.CreateDirIfNotExist(dockerPluginPath, false, 0755); err != nil {
		log.Fatal(err)
	}

	if !filedir.FileExists(dockerPluginFile) {
		err := os.WriteFile(dockerPluginFile, pluginSpecContent, 0644)
		if err != nil {
			log.Fatal(err)
		}
	}

	log.Info("Server has completed initialization")
}

func runStart(cmd *cobra.Command, args []string) {
	serverInitConfig()

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)
	signal.Notify(ch, syscall.SIGTERM)

	go func() {
		p, err := plugin.NewPlugin(adf.AppPath)
		if err != nil {
			log.Fatal(err)
		}

		h := authorization.NewHandler(p)

		log.WithFields(log.Fields{
			"version": version.Version,
		}).Info("HBM server")

		log.Info("Listening on socket file")
		log.Fatal(h.ServeUnix("hbm", 0))
	}()

	s := <-ch
	log.Infof("Processing signal '%s'", s)
}

var serverDescription = `
Launch the HBM server

`
