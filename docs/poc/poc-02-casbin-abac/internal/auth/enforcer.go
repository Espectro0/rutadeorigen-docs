package auth

import (
	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/model"
)

const modelPath = "internal/auth/model.conf"

func NewEnforcer() (*casbin.Enforcer, error) {
	m, err := model.NewModelFromFile(modelPath)
	if err != nil {
		return nil, err
	}

	return casbin.NewEnforcer(m)
}
