package auth

import "github.com/kataras/iris"

/*
   Creation Time: 2019 - Sep - 21
   Created by:  (ehsan)
   Maintainers:
      1.  Ehsan N. Moosa (E2)
   Auditor: Ehsan N. Moosa (E2)
   Copyright Ronak Software Group 2018
*/

func HasReadAccess(ctx iris.Context) bool {
	auth, ok := ctx.Values().Get(CtxAuth).(Auth)
	if !ok {
		return false
	}
	for _, p := range auth.Permissions {
		if p == Read || p == Admin {
			return true
		}
	}
	return false
}

func HasWriteAccess(ctx iris.Context) bool {
	auth, ok := ctx.Values().Get(CtxAuth).(Auth)
	if !ok {
		return false
	}
	for _, p := range auth.Permissions {
		if p == Write || p == Admin {
			return true
		}
	}
	return false
}

func HasAdminAccess(ctx iris.Context) bool {
	auth, ok := ctx.Values().Get(CtxAuth).(Auth)
	if !ok {
		return false
	}
	for _, p := range auth.Permissions {
		if p == Admin {
			return true
		}
	}
	return false
}
