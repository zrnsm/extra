// Copyright 2018 The Periph Authors. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

package hostextra

import (
	"periph.io/x/conn/v3/driver/driverreg"
	_ "periph.io/x/extra/hostextra/d2xx"
	"periph.io/x/host/v3"
)

// Init calls host.Init(), which calls driverreg.Init() and returns it as-is.
//
// The difference with host.Init() and driverreg.Init() is that hostextra.Init()
// includes more drivers, the drivers that either depend on third party
// packages or on cgo.
//
// Since host.Init() is used, all drivers in periph.io/x/host/v3 are also
// automatically loaded.
func Init() (*driverreg.State, error) {
	return host.Init()
}
