package server

import "github.com/google/wire"

var Constructor = wire.NewSet(NewRpcServer)