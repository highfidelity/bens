// Copyright 2018 High Fidelity, Inc.
//
// Distributed under the Apache License, Version 2.0.
// See the accompanying file LICENSE or http://www.apache.org/licenses/LICENSE-2.0.html

package key

import "syscall"

func mlock(data []byte) error {
	return syscall.Mlock(data)
}
