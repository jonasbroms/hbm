package server

import (
	"log/slog"
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
		slog.Error("Docker does not seem to be installed. Please check your installation.")
		os.Exit(1)
	}

	if err := filedir.CreateDirIfNotExist(dockerPluginPath, false, 0755); err != nil {
		slog.Error("Failed to create plugin directory", "error", err)
		os.Exit(1)
	}

	if !filedir.FileExists(dockerPluginFile) {
		err := os.WriteFile(dockerPluginFile, pluginSpecContent, 0644)
		if err != nil {
			slog.Error("Failed to write plugin spec file", "error", err)
			os.Exit(1)
		}
	}

	slog.Info("Server has completed initialization")
}

func runStart(cmd *cobra.Command, args []string) {
	serverInitConfig()

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)
	signal.Notify(ch, syscall.SIGTERM)

	go func() {
		p, err := plugin.NewPlugin(adf.AppPath)
		if err != nil {
			slog.Error("Failed to create plugin", "error", err)
			os.Exit(1)
		}

		h := authorization.NewHandler(p)

		slog.Info("HBM server starting",
			"version", version.Version,
			"app_path", adf.AppPath)

		slog.Info("Listening on socket",
			"socket", "unix:///run/docker/plugins/hbm.sock")

		if err := h.ServeUnix("hbm", 0); err != nil {
			slog.Error("Server failed", "error", err)
			os.Exit(1)
		}
	}()

	s := <-ch
	slog.Info("Received shutdown signal", "signal", s.String())
}

var serverDescription = `
Launch the HBM server

`
