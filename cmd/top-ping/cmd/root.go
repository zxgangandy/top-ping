package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
)

var rootCmd = &cobra.Command{
	Use:   "top-ping api server",
	Short: "top-ping api server",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

var configFile string
var config *viper.Viper

func init() {
	cobra.OnInitialize(func() {
		err := initConfig()
		if err != nil {
			panic(err)
		}
	})

	initFlags()

	rootCmd.AddCommand(serverCmd)
}

func initConfig() (err error) {
	v := viper.New()
	v.SetConfigFile(configFile)
	v.SetConfigType("yml")

	if err := v.ReadInConfig(); err != nil {
		fmt.Println(err)
		return err
	}

	config = v

	return nil
}

func initFlags() {
	rootCmd.PersistentFlags().StringVarP(&configFile, "config", "c",
		"configs/application.yml", "set config file")
}
