// Copyright 2018 High Fidelity, Inc.
//
// Distributed under the Apache License, Version 2.0.
// See the accompanying file LICENSE or http://www.apache.org/licenses/LICENSE-2.0.html

package env

import "fmt"

type EnvironmentalVariableSerializer interface {
	ToString(string, string) string
}

type ShellSerializer struct{}

func (_ ShellSerializer) ToString(name, value string) string {
	return fmt.Sprintf("export %s=\"%s\"", name, value)
}

type PowerShellSerializer struct{}

func (_ PowerShellSerializer) ToString(name, value string) string {
	return fmt.Sprintf("$env:%s = \"%s\"", name, value)
}

type CMDSerializer struct{}

func (_ CMDSerializer) ToString(name, value string) string {
	return fmt.Sprintf("set \"%s=%s\"", name, value)
}

func GetSerializer(key string) (EnvironmentalVariableSerializer, error) {
	switch key {
	case "powershell":
		return PowerShellSerializer{}, nil
	case "shell":
		return ShellSerializer{}, nil
	case "cmd":
		return CMDSerializer{}, nil
	default:
		return nil, fmt.Errorf("no serializer for %s", key)
	}
}
