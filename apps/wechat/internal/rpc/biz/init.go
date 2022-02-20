package biz

import "github.com/google/wire"

var Constructor = wire.NewSet(
	NewUserBiz,
)
