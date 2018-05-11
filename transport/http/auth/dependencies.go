package auth

import (
	"github.com/tomwright/monzoroundup/user"
	"github.com/tomwright/monzoroundup/token"
)

var (
	tokenModel token.Model
	userModel  user.Model
)

func InjectDependencies(tm token.Model, um user.Model) {
	userModel = um
	tokenModel = tm
}
