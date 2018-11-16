// Copyright 2018 High Fidelity, Inc.
//
// Distributed under the Apache License, Version 2.0.
// See the accompanying file LICENSE or http://www.apache.org/licenses/LICENSE-2.0.html

package cmd

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"

	"github.com/highfidelity/bens/cnf"
	"github.com/highfidelity/bens/env"
	"github.com/highfidelity/bens/key"
)

var serializerType string

func init() {
	rootCmd.AddCommand(environmentCmd)
	environmentCmd.PersistentFlags().StringVarP(
		&serializerType,
		"formatter", "", "shell", "choices are: shell, powershell and cmd")
}

var environmentCmd = &cobra.Command{
	Use:   "environment",
	Short: "Decrypt and display the environment",
	Run: func(cmd *cobra.Command, args []string) {
		serializer, err := env.GetSerializer(serializerType)
		if err != nil {
			log.Fatalf("couldn't load serializer: %v", err)
		}

		cipher, err := key.New(passPath, priKeyPath, pubKeyPath)
		if err != nil {
			log.Fatalf("couldn't key: %v", err)
			return
		}

		c, err := cnf.New(yamlPath, &cipher)
		environment, err := c.DecryptEnvironment()
		if err != nil {
			log.Fatalf("couldn't decrypt environment: %v", err)
		}

		for _, envVar := range environment {
			fmt.Println(serializer.ToString(envVar.Name, envVar.Value))
		}
	},
}
