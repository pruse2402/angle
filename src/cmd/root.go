/*
Copyright © 2019 NAME HERE <EMAIL ADDRESS>

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
	"log"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

// rootCmd represents the `slack` cmd (executable)
var rootCmd = &cobra.Command{
	Use:   "angle",
	Short: "APIs for angle service",
	Long: `
	APIs for APIs for angle service
	`,
}

func init() {
	// viper.SetConfigName("config")
	// viper.AddConfigPath("./config")
	// if err := viper.ReadInConfig(); err != nil {
	// 	log.Fatal("failed to read in config file")
	// }
	rootCmd.PersistentFlags().StringP(flagEnv, "e", "dev", "The environment to run in (eg. dev, test, staging, prod)")
	viper.BindPFlag(flagEnv, rootCmd.PersistentFlags().Lookup(flagEnv))
	// viper.BindEnv(flagEnv, slackEnv)

	//viper.BindEnv(flagEnv, serverEnv)

	rootCmd.PersistentFlags().StringP(flagLevel, "l", "INFO", "The minimum level of logs which are written")
	viper.BindPFlag(flagLevel, rootCmd.PersistentFlags().Lookup(flagLevel))
	// viper.BindEnv(flagLevel, slackLogLevel)

	//viper.BindEnv(flagLevel, serverLogLevel)

}

// Execute is the entry into the CLI, executing the root CMD.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
