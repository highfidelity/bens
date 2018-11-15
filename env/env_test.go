// Copyright 2018 High Fidelity, Inc.
//
// Distributed under the Apache License, Version 2.0.
// See the accompanying file LICENSE or http://www.apache.org/licenses/LICENSE-2.0.html

package env

import "testing"

func TestShellSerializer(t *testing.T) {
	expectedV := "export FOO=\"bar\""
	v := ShellSerializer{}.ToString("FOO", "bar")
	if v != expectedV {
		t.Fatalf("expected %s, got %s", expectedV, v)
	}
}

func TestPowerShellSerializer(t *testing.T) {
	expectedV := "$env:FOO = \"bar\""
	v := PowerShellSerializer{}.ToString("FOO", "bar")
	if v != expectedV {
		t.Fatalf("expected %s, got %s", expectedV, v)
	}
}

func TestCMDSerializer(t *testing.T) {
	expectedV := "set \"FOO=bar\""
	v := CMDSerializer{}.ToString("FOO", "bar")
	if v != expectedV {
		t.Fatalf("expected %s, got %s", expectedV, v)
	}
}
