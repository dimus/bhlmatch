/*
Copyright © 2020 NAME HERE <EMAIL ADDRESS>

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
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/gdower/bhlmatch"
	"github.com/spf13/cobra"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

type config struct {
	InputFile      string
	OutputFile     string
	DbHost         string
	DbUser         string
	DbPass         string
	DbName         string
	BHLnamesDir    string
	JobsNum        int
	TaxonomicMatch bool
}

var (
	cfgFile string
	opts    []bhlmatch.Option
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "bhlmatch",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {
		bhlm := bhlmatch.NewBHLmatch(opts...)
		fmt.Println(bhlm)
		bhlm.Run()
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.bhlmatch.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	// rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	home, err := homedir.Dir()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Search config in home directory with name ".bhlmatch" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".bhlmatch")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	} else {
		path := filepath.Join(home, ".bhlmatch.yaml")
		genConfig(path)
		if err := viper.ReadInConfig(); err == nil {
			fmt.Println("Default configuration needs to be changed")
		}
	}
	opts = getOpts()
}

func genConfig(path string) {
	conf := []byte(`---

# Input for CoL names and references
InputFile: /tmp/scientific_names.csv

# Output for found BHL references
OutputFile: /tmp/scientific_names_with_bhl.csv

# Postgresql host
DbHost: localhost

# Postgresql user
DbUser: postgres

# Postgresql password
DbPass:

# Postgresql database
DbName: bhlnames

# BHLnames helper directory
BHLnamesDir: /tmp/bhlnames

# Number of CPU threads
JobsNum: 4

# Match taxonomically or nomenclaturally? [Default: false]
TaxonomicMatch: false
`)

	fmt.Println("Could not find configuration file.")
	fmt.Printf("Made default config at %s, please make changes to it.\n", path)
	err := ioutil.WriteFile(path, conf, 0644)
	if err != nil {
		os.Exit(1)
	}
}

// getOpts imports data from the configuration file. These settings can be
// overriden by command line flags.
func getOpts() []bhlmatch.Option {
	var opts []bhlmatch.Option
	cfg := &config{}
	err := viper.Unmarshal(cfg)
	if err != nil {
		log.Fatal(err)
	}

	if cfg.InputFile != "" {
		opts = append(opts, bhlmatch.OptInputFile(cfg.InputFile))
	}
	if cfg.OutputFile != "" {
		opts = append(opts, bhlmatch.OptOutputFile(cfg.OutputFile))
	}
	if cfg.DbHost != "" {
		opts = append(opts, bhlmatch.OptDbHost(cfg.DbHost))
	}
	if cfg.DbUser != "" {
		opts = append(opts, bhlmatch.OptDbUser(cfg.DbUser))
	}
	if cfg.DbPass != "" {
		opts = append(opts, bhlmatch.OptDbPass(cfg.DbPass))
	}
	if cfg.DbName != "" {
		opts = append(opts, bhlmatch.OptDbName(cfg.DbName))
	}
	if cfg.BHLnamesDir != "" {
		opts = append(opts, bhlmatch.OptBHLnamesDir(cfg.BHLnamesDir))
	}
	if cfg.JobsNum != 0 {
		opts = append(opts, bhlmatch.OptJobsNum(cfg.JobsNum))
	}
	if cfg.TaxonomicMatch {
		opts = append(opts, bhlmatch.OptTaxonomicMatch(cfg.TaxonomicMatch))
	}
	return opts
}
