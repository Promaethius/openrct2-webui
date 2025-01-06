/*
Copyright © 2025 Jonathan Bryant

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/promaethius/openrct2-webui/pkg/plugin"
	"github.com/promaethius/openrct2-webui/pkg/screenshots"
	"github.com/promaethius/openrct2-webui/pkg/server"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "openrct2-webui",
	Short: "Serves giant screenshots of your park and provides a user authenticated console.",
	RunE:  runE,
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		slog.Error("an error occured", slog.String("err", err.Error()))
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().String("addr", "127.0.0.1:80", "Address to bind to for serving http traffic.")
	rootCmd.Flags().String("screenshot-directory", "/mnt/screenshots", "Screenshot directory for scanning and serving.")
	rootCmd.Flags().Uint32("screenshot-retain", 500, "Screenshots to retain in memory.")
	rootCmd.Flags().Duration("screenshot-interval", 5*time.Second, "Interval at which to generate a screenshot.")
	rootCmd.Flags().String("plugin-addr", "127.0.0.1:35711", "Address of the plugin running within the OpenRCT2 server.")
}

func runE(cmd *cobra.Command, args []string) error {
	ctx, cancel := context.WithCancel(cmd.Context())
	defer cancel()

	addr, err := cmd.Flags().GetString("addr")
	if err != nil {
		return err
	}

	screenshotDir, err := cmd.Flags().GetString("screenshot-directory")
	if err != nil {
		return err
	}

	screenshotRetain, err := cmd.Flags().GetUint32("screenshot-retain")
	if err != nil {
		return err
	}

	screenshotInterval, err := cmd.Flags().GetDuration("screenshot-interval")
	if err != nil {
		return err
	}

	pluginAddr, err := cmd.Flags().GetString("plugin-addr")
	if err != nil {
		return err
	}

	pluginClient := plugin.NewClient(pluginAddr, nil)

	screenshotManager, err := screenshots.NewManager(pluginClient, screenshotDir, screenshotInterval, screenshotRetain)
	if err != nil {
		return err
	}

	server := server.NewServer(addr, screenshotManager.GetScreenshots)

	sigs := make(chan os.Signal, 1)
	defer close(sigs)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := screenshotManager.Run(ctx); err != nil {
			slog.Error("an error occured within screenshot manager", slog.String("error", err.Error()))
		}
	}()

	go func() {
		if err := server.Run(); err != nil {
			slog.Error("an error occured within server", slog.String("error", err.Error()))
		}
	}()

	<-sigs
	slog.Info("shutdown signal detected")

	server.Shutdown(ctx)

	return nil
}
