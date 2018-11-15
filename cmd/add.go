// Copyright 2018 High Fidelity, Inc.
//
// Distributed under the Apache License, Version 2.0.
// See the accompanying file LICENSE or http://www.apache.org/licenses/LICENSE-2.0.html

package cmd

import (
	"errors"
	"log"

	"github.com/spf13/cobra"

	"github.com/highfidelity/bens/cnf"
	"github.com/highfidelity/bens/key"
)

func init() {
	rootCmd.AddCommand(addCmd)
}

var addCmd = &cobra.Command{
	Use:   "add ENV_NAME ENV_VALUE",
	Short: "Add an encrypted environment variable to the environment",
	Args: func(cmd *cobra.Command, args []string) error {
		l := len(args)
		if l == 0 {
			return errors.New("supply the ENV_NAME and ENV_VALUE arguments")
		} else if l == 1 {
			return errors.New("supply the ENV_VALUE argument")
		} else if l > 2 {
			return errors.New("too many arguments supplied")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		envName := args[0]
		envValue := args[1]

		cipher, err := key.New(passPath, priKeyPath, pubKeyPath)
		if err != nil {
			log.Fatalf("couldn't load key: %v", err)
			return
		}

		cnf, err := cnf.New(yamlPath, cipher)
		if err != nil {
			log.Fatalf("couldn't load yaml %s: %v", yamlPath, err)
		}

		if err = cnf.Add(envName, envValue); err != nil {
			log.Fatalf("couldn't add environment variable: %v", err)
		}
		if err = cnf.Save(yamlPath); err != nil {
			log.Fatalf("couldn't save yaml to %s: %v", yamlPath, err)
		}
	},
}
