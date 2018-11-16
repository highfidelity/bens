// Copyright 2018 High Fidelity, Inc.
//
// Distributed under the Apache License, Version 2.0.
// See the accompanying file LICENSE or http://www.apache.org/licenses/LICENSE-2.0.html

package cnf

import (
	"fmt"
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

type yamlEnvVar struct {
	Name           string
	EncryptedValue string `yaml:"encryptedValue"`
}

type yamlRoot struct {
	Version     int
	Environment []yamlEnvVar
}

type cipher interface {
	Decrypt(string) (string, error)
	Encrypt(string) (string, error)
}

type Cnf struct {
	root   yamlRoot
	cipher cipher
}

type EnvVar struct {
	Name  string
	Value string
}

func (c *Cnf) Add(name, value string) error {
	cipherText, err := c.cipher.Encrypt(value)
	if err != nil {
		return fmt.Errorf("couldn't encrypt %s: %v", value, err)
	}
	v := yamlEnvVar{Name: name, EncryptedValue: cipherText}
	c.root.Environment = append(c.root.Environment, v)
	return nil
}

func (c *Cnf) Save(yamlPath string) error {
	out, err := yaml.Marshal(&c.root)
	if err != nil {
		return err
	}
	if err := ioutil.WriteFile(yamlPath, out, 0644); err != nil {
		return fmt.Errorf("couldn't save to %s: %v", yamlPath, err)
	}
	return nil
}

func (c *Cnf) DecryptEnvironment() ([]EnvVar, error) {
	env := make([]EnvVar, 0)
	for _, envVar := range c.root.Environment {
		val, err := c.cipher.Decrypt(envVar.EncryptedValue)
		if err != nil {
			return nil, fmt.Errorf("couldn't decrypt %s: %v", envVar.Name, err)
		}
		env = append(env, EnvVar{Name: envVar.Name, Value: val})
	}
	return env, nil
}

func New(yamlPath string, e cipher) (Cnf, error) {
	c := Cnf{cipher: e}
	y, err := ioutil.ReadFile(yamlPath)
	if err != nil {
		return c, fmt.Errorf("couldn't open yml: %v", err)
	}
	err = yaml.Unmarshal(y, &c.root)
	if err != nil {
		return c, fmt.Errorf("couldn't unmarshall yaml: %v", err)
	}
	return c, nil
}
