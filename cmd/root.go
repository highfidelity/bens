// Copyright 2018 High Fidelity, Inc.
//
// Distributed under the Apache License, Version 2.0.
// See the accompanying file LICENSE or http://www.apache.org/licenses/LICENSE-2.0.html

package cmd

import (
	"log"

	"github.com/spf13/cobra"
)

var yamlPath, passPath, priKeyPath, pubKeyPath string

func init() {
	rootCmd.PersistentFlags().StringVarP(&yamlPath, "config-file", "c", "bens.yml", "configuration file")
	rootCmd.PersistentFlags().StringVarP(&passPath, "pass-file", "p", "pass.txt", "pass file")
	rootCmd.PersistentFlags().StringVarP(&priKeyPath, "private-key-file", "", "pri.key", "private key file")
	rootCmd.PersistentFlags().StringVarP(&pubKeyPath, "public-key-file", "", "pub.key", "public key file")
}

var rootCmd = &cobra.Command{
	Use:   "bens",
	Short: "build environment security",
	Run:   func(cmd *cobra.Command, args []string) {},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatalf("error: %v", err)
	}
}
