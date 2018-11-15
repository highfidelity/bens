// Copyright 2018 High Fidelity, Inc.
//
// Distributed under the Apache License, Version 2.0.
// See the accompanying file LICENSE or http://www.apache.org/licenses/LICENSE-2.0.html

package cnf

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

var yamlFileContent = []byte(`
version: 1
environment:
  - name: FOO
    encryptedValue: foo
`)

func writeTestFiles() (string, string, error) {
	dir, err := ioutil.TempDir("", "cnf_test")
	if err != nil {
		return "", "", err
	}

	yamlFilePath := filepath.Join(dir, "bens.yml")
	if ioutil.WriteFile(yamlFilePath, yamlFileContent, 0644) != nil {
		os.RemoveAll(dir)
		return "", "", err
	}

	return dir, yamlFilePath, nil
}

type MockCipher struct {
}

func (m MockCipher) Decrypt(cipherText string) (string, error) {
	return cipherText + "bar", nil
}

func (m MockCipher) Encrypt(plainText string) (string, error) {
	return "foo" + plainText, nil
}

func TestDecryptEnvironment(t *testing.T) {
	dir, yamlFilePath, err := writeTestFiles()
	if err != nil {
		t.Fatalf("couldn't write test files: %v", err)
	}
	defer os.RemoveAll(dir)

	mockCipher := MockCipher{}

	cnf, err := New(yamlFilePath, mockCipher)
	if err != nil {
		t.Fatalf("couldn't create Cnf with New: %v", err)
	}
	environment, err := cnf.DecryptEnvironment()
	if err != nil {
		t.Fatalf("couldn't decrypt environment: %v", err)
	}
	if len(environment) != 1 {
		t.Fatalf("environment should contain one variable; contains %d", len(environment))
	}
	if environment[0].Name != "FOO" {
		t.Fatalf("environment name should be FOO; is: %v", environment[0].Name)
	}
	if environment[0].Value != "foobar" {
		t.Fatalf("environment value should be foobar; is: %v", environment[0].Value)
	}
}

func TestAdd(t *testing.T) {
	dir, yamlFilePath, err := writeTestFiles()
	if err != nil {
		t.Fatalf("couldn't write test files: %v", err)
	}
	defer os.RemoveAll(dir)

	cnf, err := New(yamlFilePath, MockCipher{})
	if err != nil {
		t.Fatalf("couldn't create Cnf with New: %v", err)
	}

	if err := cnf.Add("BAR", "bar"); err != nil {
		t.Fatalf("couldn't Add environmental variable: %v", err)
	}

	environment, err := cnf.DecryptEnvironment()
	if err != nil {
		t.Fatalf("couldn't get environment after Add: %v", err)
	}
	if len(environment) != 2 {
		t.Fatalf("environment should have two values after Add")
	}
	v := environment[1]
	if v.Name != "BAR" {
		t.Fatalf("new environment variable should be named BAR: was %s", v.Name)
	}
	if v.Value != "foobarbar" {
		t.Fatalf("new environment variable should be foobarbar: was %s", v.Value)
	}
}

func TestSave(t *testing.T) {
	dir, yamlFilePath, err := writeTestFiles()
	if err != nil {
		t.Fatalf("couldn't write test files: %v", err)
	}
	defer os.RemoveAll(dir)

	cnf, err := New(yamlFilePath, MockCipher{})
	if err != nil {
		t.Fatalf("couldn't create Cnf with New: %v", err)
	}

	cnf.root.Version = 12345
	cnf.Save(yamlFilePath)

	cnf, err = New(yamlFilePath, MockCipher{})
	if err != nil {
		t.Fatalf("reloading Cnf with New failed: %v", err)
	}

	if cnf.root.Version != 12345 {
		t.Fatal("changes to cnf didn't persist")
	}
}
