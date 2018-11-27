// Copyright 2018 High Fidelity, Inc.
//
// Distributed under the Apache License, Version 2.0.
// See the accompanying file LICENSE or http://www.apache.org/licenses/LICENSE-2.0.html

package cmd

import (
	"fmt"
	"log"
	"os"
	"syscall"

	"github.com/spf13/cobra"
	"golang.org/x/crypto/ssh/terminal"

	"github.com/highfidelity/bens/cnf"
	"github.com/highfidelity/bens/env"
	"github.com/highfidelity/bens/key"
)

func readPassFromTerm() ([]byte, error) {
	fmt.Printf("password: ")
	pass, err := terminal.ReadPassword(int(syscall.Stdin))
	if err != nil {
		return nil, err
	}
	fmt.Println("")
	return pass, nil
}

var serializerType string
var shouldAskPass bool

func init() {
	rootCmd.AddCommand(environmentCmd)
	environmentCmd.PersistentFlags().StringVarP(
		&serializerType,
		"formatter", "", "shell", "choices are: shell, powershell and cmd")
	environmentCmd.PersistentFlags().BoolVarP(
		&shouldAskPass, "ask-pass", "", false, "ask for pass")
}

var environmentCmd = &cobra.Command{
	Use:   "environment",
	Short: "Decrypt and display the environment",
	Run: func(cmd *cobra.Command, args []string) {
		serializer, err := env.GetSerializer(serializerType)
		if err != nil {
			log.Fatalf("couldn't load serializer: %v", err)
		}

		var cipher key.Key
		if shouldAskPass {
			pass, err := readPassFromTerm()
			if err != nil {
				log.Fatalf("couldn't read pass from terminal: %v", err)
			}
			cipher, err = key.NewWithPass(pass, priKeyPath, pubKeyPath)
		} else {
			pass := os.Getenv("BENS_PASS")
			if pass != "" {
				cipher, err = key.NewWithPass([]byte(pass), priKeyPath, pubKeyPath)
			} else {
				cipher, err = key.New(passPath, priKeyPath, pubKeyPath)
			}
		}
		if err != nil {
			log.Fatalf("couldn't read key: %v", err)
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
