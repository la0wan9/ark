package main

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/fsnotify/fsnotify"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/automaxprocs/maxprocs"

	"github.com/la0wan9/ark/cmd"
)

var version string

func main() {
	checkError(initViper())
	checkError(initAutomaxprocs())
	cobra.EnableCommandSorting = false
	rootCmd := &cobra.Command{
		Use:     filepath.Base(os.Args[0]),
		Version: version,
	}
	rootCmd.AddCommand(cmd.NewServerCmd())
	rootCmd.AddCommand(cmd.NewSpiderCmd())
	checkError(rootCmd.Execute())
}

func checkError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func initAutomaxprocs() error {
	_, err := maxprocs.Set(maxprocs.Logger(nil))
	return err
}

func initViper() error {
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AddConfigPath(".")
	viper.AddConfigPath("..")
	if err := viper.ReadInConfig(); err != nil {
		return err
	}
	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		log.Infof("config file changed: %s", e.Name)
	})
	return nil
}
